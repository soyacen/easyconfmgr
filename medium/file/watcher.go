package easyconfmgrfile

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

func (w *Watcher) Watch() error {
	w.log.Info("start watch file:", w.filename)
	err := w.watcher.Add(w.filename)
	if err != nil {
		return fmt.Errorf("failed add watcher file, %w", err)
	}
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					w.log.Info("stop watch file:", w.filename)
					return
				}
				w.log.Infof("an event [%s] occurs", event.String())
				if event.Op&fsnotify.Write == fsnotify.Write {
					data, err := ioutil.ReadFile(w.filename)
					if err != nil {
						w.log.Error(fmt.Errorf("failed read file %s, %w", w.filename, err))
						return
					}
					w.events <- easyconfmgr.NewEvent(event, data)
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					w.log.Error(fmt.Errorf("failed watch, %w", err))
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	return nil
}

func (w *Watcher) Events() <-chan *easyconfmgr.Event {
	return w.events
}

func (w *Watcher) Stop() error {
	return w.stop()
}

func (w *Watcher) stop() error {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if !w.isStopped {
		w.log.Info("stop watching...")
		w.isStopped = true
		close(w.events)
		return w.watcher.Close()
	}
	return nil
}

func NewWatcher(filename string, log easyconfmgr.Logger) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed new watcher, %w", err)
	}
	w := &Watcher{watcher: watcher, filename: filename, log: log, events: make(chan *easyconfmgr.Event)}
	return w, nil
}
