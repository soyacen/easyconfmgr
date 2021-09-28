package nacos

import (
	"errors"
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
	watchOnce sync.Once
	stopOnce  sync.Once
	log       easyconfmgr.Logger
}

func (watcher *Watcher) Watch() error {
	err := errors.New("watcher had started watch")
	watcher.watchOnce.Do(func() {
		err = watcher.watch()
	})
	return err
}

func (watcher *Watcher) watch() error {
	watcher.log.Infof("start listen config, dataId: %s, group: %s", watcher.dataId, watcher.group)
	err := watcher.client.ListenConfig(vo.ConfigParam{
		DataId: watcher.dataId,
		Group:  watcher.group,
		OnChange: func(namespace, group, dataId, data string) {
			event := &EventDescription{namespace: namespace, group: group, dataId: dataId}
			watcher.events <- easyconfmgr.NewEvent(event, []byte(data))
		},
	})
	if err != nil {
		return fmt.Errorf("failed to listen config, DataId: %s, Group: %s , %w", watcher.dataId, watcher.group, err)
	}
	return nil
}

func (watcher *Watcher) Events() <-chan *easyconfmgr.Event {
	return watcher.events
}

func (watcher *Watcher) Stop() error {
	err := errors.New("watcher had stopped")
	watcher.stopOnce.Do(func() {
		err = watcher.stop()
	})
	return err
}

func (watcher *Watcher) stop() error {
	watcher.log.Info("stop watching nacos config...")
	close(watcher.events)
	err := watcher.client.CancelListenConfig(vo.ConfigParam{
		DataId: watcher.dataId,
		Group:  watcher.group,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel listen config, %w", err)
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
	watcher := &Watcher{
		client:    client,
		group:     group,
		dataId:    dataId,
		events:    make(chan *easyconfmgr.Event),
		watchOnce: sync.Once{},
		stopOnce:  sync.Once{},
		log:       nil,
	}
	for _, opt := range opts {
		opt(watcher)
	}
	return watcher
}
