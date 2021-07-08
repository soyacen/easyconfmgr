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
	l.log.Infof("get config DataId: %s, Group: %s", l.dataId, l.group)
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
	log easyconfmgr.Logger,
) easyconfmgr.Loader {
	return &Loader{client: client, group: group, dataId: dataId, contentType: contentType, log: log}
}
