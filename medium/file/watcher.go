package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/soyacen/easyconfmgr"
)

type Watcher struct {
	filename  string
	watcher   *fsnotify.Watcher
	events    chan *easyconfmgr.Event
	isStopped bool
	mu        sync.RWMutex
	log       easyconfmgr.Logger
}

func (watcher *Watcher) Watch() error {
	watcher.log.Info("start watch file:", watcher.filename)
	err := watcher.watcher.Add(watcher.filename)
	if err != nil {
		return fmt.Errorf("failed add watcher file, %w", err)
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.watcher.Events:
				if !ok {
					watcher.log.Info("stop watch file:", watcher.filename)
					return
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
			case err, ok := <-watcher.watcher.Errors:
				if !ok {
					watcher.log.Error(fmt.Errorf("failed watch, %w", err))
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	return nil
}

func (watcher *Watcher) Events() <-chan *easyconfmgr.Event {
	return watcher.events
}

func (watcher *Watcher) Stop() error {
	return watcher.stop()
}

func (watcher *Watcher) stop() error {
	watcher.mu.RLock()
	defer watcher.mu.RUnlock()
	if !watcher.isStopped {
		watcher.log.Info("stop watching...")
		watcher.isStopped = true
		close(watcher.events)
		return watcher.watcher.Close()
	}
	return nil
}

type WatcherOption func(watcher *Watcher)

func WithLogger(log easyconfmgr.Logger) LoaderOption {
	return func(loader *Loader) {
		loader.log = log
	}
}

func NewWatcher(filename string, opts ...WatcherOption) (*Watcher, error) {
	w := &Watcher{filename: filename, events: make(chan *easyconfmgr.Event), log: easyconfmgr.DiscardLogger}
	for _, opt := range opts {
		opt(w)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed new watcher, %w", err)
	}
	w.watcher = watcher
	return w, nil
}
