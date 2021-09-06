package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/soyacen/easyconfmgr"
)

// nacos配置路径定义
// nacos://address:port?contentType=yaml&namespace=ns&group=g&dataId=d

type Loader struct {
	client      config_client.IConfigClient
	group       string
	dataId      string
	contentType string
	data        []byte
	log         easyconfmgr.Logger
}

func (loader *Loader) ContentType() string {
	return loader.contentType
}

func (loader *Loader) Load() error {
	loader.log.Infof("get config DataId: %s, Group: %s", loader.dataId, loader.group)
	content, err := loader.client.GetConfig(vo.ConfigParam{
		DataId: loader.dataId,
		Group:  loader.group,
	})
	if err != nil {
		return err
	}
	loader.data = []byte(content)
	return nil
}

func (loader *Loader) RawData() []byte {
	return loader.data
}

type LoaderOption func(loader *Loader)

func Logger(log easyconfmgr.Logger) LoaderOption {
	return func(loader *Loader) {
		loader.log = log
	}
}

func NewLoader(
	client config_client.IConfigClient,
	group string,
	dataId string,
	contentType string,
	opts ...LoaderOption,
) easyconfmgr.Loader {
	loader := &Loader{client: client, group: group, dataId: dataId, contentType: contentType, log: easyconfmgr.DiscardLogger}
	for _, opt := range opts {
		opt(loader)
	}
	return loader
}
