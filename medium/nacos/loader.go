package easyconfmgrnacos

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

func (l *Loader) ContentType() string {
	return l.contentType
}

func (l *Loader) Load() error {
	content, err := l.client.GetConfig(vo.ConfigParam{
		DataId: l.dataId,
		Group:  l.group,
	})
	if err != nil {
		return err
	}
	l.data = []byte(content)
	return nil
}

func (l *Loader) RawData() []byte {
	return l.data
}

func NewLoader(
	client config_client.IConfigClient,
	group string,
	dataId string,
	contentType string,
) easyconfmgr.Loader {
	return &Loader{client: client, group: group, dataId: dataId, contentType: contentType}
}
