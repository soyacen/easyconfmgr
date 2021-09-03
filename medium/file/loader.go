package mediumfile

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/soyacen/goutils/stringutils"

	"github.com/soyacen/easyconfmgr"
)

type Loader struct {
	filename    string
	contentType string
	rawData     []byte
	log         easyconfmgr.Logger
}

func (loader *Loader) ContentType() string {
	return loader.contentType
}

func (loader *Loader) Load() error {
	loader.log.Info("reading file:", loader.filename)
	rawData, err := ioutil.ReadFile(loader.filename)
	if err != nil {
		return err
	}
	loader.log.Debug("file content:", string(rawData))
	loader.rawData = rawData
	return nil
}

func (loader *Loader) RawData() []byte {
	return loader.rawData
}

type LoaderOption func(loader *Loader)

func ContentType(contentType string) LoaderOption {
	return func(loader *Loader) {
		loader.contentType = contentType
	}
}

func Logger(log easyconfmgr.Logger) LoaderOption {
	return func(loader *Loader) {
		loader.log = log
	}
}

func NewLoader(filename string, opts ...LoaderOption) easyconfmgr.Loader {
	loader := &Loader{filename: filename, log: easyconfmgr.DiscardLogger}
	for _, opt := range opts {
		opt(loader)
	}
	if stringutils.IsBlank(loader.contentType) {
		loader.contentType = filepath.Ext(filename)
		if strings.HasPrefix(loader.contentType, ".") {
			loader.contentType = loader.contentType[1:]
		}
	}
	return loader
}
