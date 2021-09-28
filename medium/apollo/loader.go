package apollo

import (
	"context"
	"fmt"
	"time"

	"github.com/soyacen/easyhttp"
	"github.com/soyacen/goutils/fileutils"
	"github.com/soyacen/goutils/stringutils"

	"github.com/soyacen/easyconfmgr"
)

type Loader struct {
	scheme        string
	host          string
	port          int
	appID         string
	cluster       string
	namespaceName string
	secret        string
	timeout       time.Duration
	client        *easyhttp.Client
	rawData       []byte
	log           easyconfmgr.Logger
}

func (loader *Loader) ContentType() string {
	contentType := fileutils.ExtName(loader.namespaceName)
	if stringutils.IsBlank(contentType) {
		return "properties"
	}
	return contentType
}

func (loader *Loader) Load() error {
	content, err := loader.getConfigFromApollo()
	if err != nil {
		return err
	}
	loader.log.Debug("config content:", content)
	loader.rawData = []byte(content)
	return nil
}

func (loader *Loader) getConfigFromApollo() (string, error) {
	uri := fmt.Sprintf("%s://%s:%d/configs/%s/%s/%s", loader.scheme, loader.host, loader.port, loader.appID, loader.cluster, loader.namespaceName)
	loader.log.Info("reading apollo config:", uri)
	ctx := context.Background()
	if loader.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, loader.timeout)
		defer cancel()
	}
	return getConfigContent(ctx, uri, loader.appID, loader.secret, loader.ContentType(), loader.client)
}

func (loader *Loader) RawData() []byte {
	return loader.rawData
}

type LoaderOption func(loader *Loader)

func Scheme(scheme string) LoaderOption {
	return func(loader *Loader) {
		loader.scheme = scheme
	}
}

func Secret(secret string) LoaderOption {
	return func(loader *Loader) {
		loader.secret = secret
	}
}

func Timeout(timeout time.Duration) LoaderOption {
	return func(loader *Loader) {
		loader.timeout = timeout
	}
}

func Logger(log easyconfmgr.Logger) LoaderOption {
	return func(loader *Loader) {
		loader.log = log
	}
}

func NewLoader(host string, port int, appID string, cluster string, namespaceName string, opts ...LoaderOption) easyconfmgr.Loader {
	loader := &Loader{
		scheme:        "http",
		host:          host,
		port:          port,
		appID:         appID,
		cluster:       cluster,
		namespaceName: namespaceName,
		log:           easyconfmgr.DiscardLogger,
	}
	for _, opt := range opts {
		opt(loader)
	}

	client := easyhttp.NewClient(easyhttp.WithTimeout(loader.timeout))
	loader.client = client

	return loader
}
