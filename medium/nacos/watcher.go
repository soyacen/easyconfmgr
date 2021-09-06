package nacos

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

func (watcher *Watcher) Watch() error {
	watcher.log.Infof("start listen config, dataId: %s, group: %s", watcher.dataId, watcher.group)
	return watcher.client.ListenConfig(vo.ConfigParam{
		DataId: watcher.dataId,
		Group:  watcher.group,
		OnChange: func(namespace, group, dataId, data string) {
			event := &EventDescription{namespace: namespace, group: group, dataId: dataId}
			watcher.events <- easyconfmgr.NewEvent(event, []byte(data))
		},
	})
}

func (watcher *Watcher) Events() <-chan *easyconfmgr.Event {
	return watcher.events
}

func (watcher *Watcher) Stop() error {
	return watcher.stop()
}

func (watcher *Watcher) stop() error {
	watcher.mu.Lock()
	defer watcher.mu.Unlock()
	if !watcher.isStopped {
		watcher.log.Info("stop watching...")
		watcher.isStopped = true
		close(watcher.events)
		return watcher.client.CancelListenConfig(vo.ConfigParam{
			DataId: watcher.dataId,
			Group:  watcher.group,
		})
	}
	return nil
}

type WatcherOption func(watcher *Watcher)

func WithLogger(log easyconfmgr.Logger) LoaderOption {
	return func(loader *Loader) {
		loader.log = log
	}
}

func NewWatcher(
	client config_client.IConfigClient,
	group string,
	dataId string,
	opts ...WatcherOption,
) *Watcher {
	watcher := &Watcher{client: client, group: group, dataId: dataId, events: make(chan *easyconfmgr.Event)}
	for _, opt := range opts {
		opt(watcher)
	}
	return watcher
}
