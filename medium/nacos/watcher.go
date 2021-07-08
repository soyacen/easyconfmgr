package easyconfmgrnacos

import (
	"fmt"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/soyacen/easyconfmgr"
)

type EventDescription struct {
	namespace string
	group     string
	dataId    string
}

func (ed *EventDescription) String() string {
	return fmt.Sprintf("%s.%s.%s", ed.namespace, ed.group, ed.dataId)
}

type Watcher struct {
	client    config_client.IConfigClient
	group     string
	dataId    string
	events    chan *easyconfmgr.Event
	isStopped bool
	mu        sync.Mutex
	log       easyconfmgr.Logger
}

func (w *Watcher) Watch() error {
	w.log.Infof("start listen config, dataId: %s, group: %s", w.dataId, w.group)
	return w.client.ListenConfig(vo.ConfigParam{
		DataId: w.dataId,
		Group:  w.group,
		OnChange: func(namespace, group, dataId, data string) {
			event := &EventDescription{namespace: namespace, group: group, dataId: dataId}
			w.events <- easyconfmgr.NewEvent(event, []byte(data))
		},
	})
}

func (w *Watcher) Events() <-chan *easyconfmgr.Event {
	return w.events
}

func (w *Watcher) Stop() error {
	return w.stop()
}

func (w *Watcher) stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.isStopped {
		w.log.Info("stop watching...")
		w.isStopped = true
		close(w.events)
		return w.client.CancelListenConfig(vo.ConfigParam{
			DataId: w.dataId,
			Group:  w.group,
		})
	}
	return nil
}

func NewWatcher(
	client config_client.IConfigClient,
	group string,
	dataId string,
	log easyconfmgr.Logger,
) *Watcher {
	return &Watcher{client: client, group: group, dataId: dataId, log: log, events: make(chan *easyconfmgr.Event)}
}
