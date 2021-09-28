package file

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/soyacen/easyconfmgr"
)

type Watcher struct {
	filename  string
	fsWatcher *fsnotify.Watcher
	events    chan *easyconfmgr.Event
	watchOnce sync.Once
	stopCtx   context.Context
	stopFunc  context.CancelFunc
	stopOnce  sync.Once
	log       easyconfmgr.Logger
}

func (watcher *Watcher) Watch() error {
	err := errors.New("watcher had started watch")
	watcher.watchOnce.Do(func() {
		err = watcher.addFileWatcher()
		go watcher.watch()
	})
	return err
}

func (watcher *Watcher) addFileWatcher() error {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed new watcher, %w", err)
	}
	watcher.fsWatcher = fsWatcher

	watcher.log.Info("start watch file:", watcher.filename)
	err = watcher.fsWatcher.Add(watcher.filename)
	if err != nil {
		return fmt.Errorf("failed add watcher file, %w", err)
	}
	return nil
}

func (watcher *Watcher) watch() {
	defer watcher.fsWatcher.Close()
	for {
		select {
		case <-watcher.stopCtx.Done():
			watcher.log.Info("stop watch file:", watcher.filename)
			return
		case event, ok := <-watcher.fsWatcher.Events:
			if !ok {
				watcher.log.Info("stop watch file:", watcher.filename)
				continue
			}
			watcher.log.Infof("an event [%s] occurs", event.String())
			if event.Op&fsnotify.Write == fsnotify.Write {
				data, err := ioutil.ReadFile(watcher.filename)
				if err != nil {
					watcher.log.Error(fmt.Errorf("failed read file %s, %w", watcher.filename, err))
					return
				}
				watcher.events <- easyconfmgr.NewEvent(event, data)
			}
		case err, ok := <-watcher.fsWatcher.Errors:
			if !ok {
				watcher.log.Info("stop watch file:", watcher.filename)
				continue
			}
			watcher.log.Error(fmt.Errorf("failed watch, %w", err))
		}
	}
	return
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
	watcher.log.Info("stop watching config...")
	watcher.stopFunc()
	return nil
}

type WatcherOption func(watcher *Watcher)

func WithLogger(log easyconfmgr.Logger) LoaderOption {
	return func(loader *Loader) {
		loader.log = log
	}
}

func NewWatcher(filename string, opts ...WatcherOption) *Watcher {
	w := &Watcher{
		filename:  filename,
		events:    make(chan *easyconfmgr.Event),
		watchOnce: sync.Once{},
		stopOnce:  sync.Once{},
		log:       easyconfmgr.DiscardLogger,
	}
	for _, opt := range opts {
		opt(w)
	}

	w.stopCtx, w.stopFunc = context.WithCancel(context.Background())

	return w
}
