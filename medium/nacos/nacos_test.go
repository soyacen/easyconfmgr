package nacos_test

import (
	"log"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/assert"

	"github.com/soyacen/easyconfmgr"
	"github.com/soyacen/easyconfmgr/medium/nacos"
)

var namespaceID = "public"
var group = "test"
var dataID = "demo"
var confType = "yaml"
var confContent = `bool: true
int: 10
int_32: -200
int_64: -3000
u_int: 133
u_int_32: 413
u_int_64: 564
float_64: 1.3
time: 2021-07-07T17:16:12.361234+08:00
duration: 1s`

var client config_client.IConfigClient

func TestMain(m *testing.M) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}
	cc := constant.NewClientConfig()
	var err error
	client, err = clients.NewConfigClient(vo.NacosClientParam{ClientConfig: cc, ServerConfigs: sc})
	if err != nil {
		log.Fatalln(err)
	}

	ok, err := client.PublishConfig(vo.ConfigParam{Group: group, DataId: dataID, Content: confContent})
	if err != nil {
		log.Fatalln(err)
	}

	if !ok {
		log.Fatalln(ok)
	}

	m.Run()

	ok, err = client.DeleteConfig(vo.ConfigParam{Group: group, DataId: dataID})
	if err != nil {
		log.Fatalln(err)
	}

	if !ok {
		log.Fatalln(ok)
	}
}

func TestLoader(t *testing.T) {
	loader := nacos.NewLoader(client, group, dataID, confType)
	contentType := loader.ContentType()
	assert.Equal(t, confType, contentType, "content type not match")

	err := loader.Load()
	assert.Nil(t, err, "failed load file")

	data := loader.RawData()
	assert.Equal(t, confContent, string(data), "config content not equal")
}

func TestWatcher(t *testing.T) {
	watcher := nacos.NewWatcher(client, group, dataID)
	err := watcher.Watch()
	assert.Nil(t, err, "failed watch file")

	go func() {
		time.Sleep(3 * time.Second)
		client, err := clients.NewConfigClient(vo.NacosClientParam{ClientConfig: constant.NewClientConfig(), ServerConfigs: []constant.ServerConfig{*constant.NewServerConfig("127.0.0.1", 8848)}})
		if err != nil {
			log.Fatalln(err)
		}
		content := confContent + "\nfloat_32: 0.3"
		ok, err := client.PublishConfig(vo.ConfigParam{Group: group, DataId: dataID, Content: content})
		assert.Nil(t, err, "failed open file")
		assert.True(t, ok, "publish config not ok")

		time.Sleep(3 * time.Second)
		err = watcher.Stop()
	}()
	events := watcher.Events()
	var e *easyconfmgr.Event
	for event := range events {
		e = event
	}
	assert.Equal(t, confContent+"\nfloat_32: 0.3", string(e.Data()))
}
