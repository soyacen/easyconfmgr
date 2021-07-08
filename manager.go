package easyconfmgr

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
)

// Manager config manager, support mutil load
type Manager struct {
	// loaders
	loaders     []Loader
	parsers     []Parser
	watchers    []Watcher
	valuer      Valuer
	mu          sync.RWMutex
	eventCh     chan *Event
	isStopWatch bool
}

type Option func(mgr *Manager)

func WithLoader(loader ...Loader) Option {
	return func(mgr *Manager) {
		mgr.loaders = append(mgr.loaders, loader...)
	}
}

func WithParser(parser ...Parser) Option {
	return func(mgr *Manager) {
		mgr.parsers = append(mgr.parsers, parser...)
	}
}

func WithWatcher(watcher ...Watcher) Option {
	return func(mgr *Manager) {
		mgr.watchers = append(mgr.watchers, watcher...)
	}
}

func WithValuer(valuer Valuer) Option {
	return func(mgr *Manager) {
		mgr.valuer = valuer
	}
}

func NewManager(opts ...Option) *Manager {
	m := &Manager{
		eventCh: make(chan *Event),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Manager) ReadConfig() error {
	loaders := m.loaders
	configs := make([]map[string]interface{}, 0, len(loaders))
	for _, loader := range loaders {
		if err := loader.Load(); err != nil {
			return err
		}
		parser, err := m.getParser(loader.ContentType())
		if err != nil {
			return err
		}
		if err := parser.Parse(loader.RawData()); err != nil {
			return err
		}
		configs = append(configs, parser.ConfigMap())
	}
	m.valuer.AddConfig(configs...)
	return nil
}

func (m *Manager) AllConfigs() map[string]interface{} {
	return m.valuer.AllConfigs()
}

func (m *Manager) getParser(contentType string) (Parser, error) {
	for _, parser := range m.parsers {
		if parser.Support(contentType) {
			return parser, nil
		}
	}
	return nil, fmt.Errorf("not found %s parser", contentType)
}

func (m *Manager) StartWatch() error {
	var err error
	for _, watcher := range m.watchers {
		if e := watcher.Watch(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	return err
}

func (m *Manager) Events() <-chan *Event {
	for _, watcher := range m.watchers {
		go func(events <-chan *Event) {
			for event := range events {
				m.eventCh <- event
			}
		}(watcher.Events())
	}
	return m.eventCh
}

func (m *Manager) StopWatch() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var err error
	if !m.isStopWatch {
		m.isStopWatch = true
		for _, watcher := range m.watchers {
			if e := watcher.Stop(); e != nil {
				err = multierror.Append(err, e)
			}
		}
		close(m.eventCh)
	}
	return err
}

// =================== Valuer ===================

func (m *Manager) Get(key string) (interface{}, error) {
	return m.valuer.Get(key)
}

func (m *Manager) UnmarshalKey(key string, rawVal interface{}) error {
	return m.valuer.UnmarshalKey(key, rawVal)
}

func (m *Manager) Unmarshal(rawVal interface{}) error {
	return m.valuer.Unmarshal(rawVal)
}

// =================== FloatValuer ===================

func (m *Manager) GetFloat32(key string) (float32, error) {
	return m.valuer.GetFloat32(key)
}

func (m *Manager) GetFloat64(key string) (float64, error) {
	return m.valuer.GetFloat64(key)
}

// =================== IntValuer ===================

func (m *Manager) GetInt(key string) (int, error) {
	return m.valuer.GetInt(key)
}

func (m *Manager) GetInt8(key string) (int8, error) {
	return m.valuer.GetInt8(key)
}

func (m *Manager) GetInt16(key string) (int16, error) {
	return m.valuer.GetInt16(key)
}

func (m *Manager) GetInt32(key string) (int32, error) {
	return m.valuer.GetInt32(key)
}

func (m *Manager) GetInt64(key string) (int64, error) {
	return m.valuer.GetInt64(key)
}

func (m *Manager) GetUint(key string) (uint, error) {
	return m.valuer.GetUint(key)
}

func (m *Manager) GetUint8(key string) (uint8, error) {
	return m.valuer.GetUint8(key)
}

func (m *Manager) GetUint16(key string) (uint16, error) {
	return m.valuer.GetUint16(key)
}

func (m *Manager) GetUint32(key string) (uint32, error) {
	return m.valuer.GetUint32(key)
}

func (m *Manager) GetUint64(key string) (uint64, error) {
	return m.valuer.GetUint64(key)
}

// =================== BoolValuer ===================

func (m *Manager) GetBool(key string) (bool, error) {
	return m.valuer.GetBool(key)
}

// =================== StringValuer ===================

func (m *Manager) GetString(key string) (string, error) {
	return m.valuer.GetString(key)
}

// =================== TimeValuer ===================

func (m *Manager) GetTime(key string) (time.Time, error) {
	return m.valuer.GetTime(key)
}

func (m *Manager) GetDuration(key string) (time.Duration, error) {
	return m.valuer.GetDuration(key)
}

// =================== SliceValuer ===================

func (m *Manager) GetSlice(key string) ([]interface{}, error) {
	return m.valuer.GetSlice(key)
}

func (m *Manager) GetIntSlice(key string) ([]int, error) {
	return m.valuer.GetIntSlice(key)
}

func (m *Manager) GetBoolSlice(key string) ([]bool, error) {
	return m.valuer.GetBoolSlice(key)
}

func (m *Manager) GetStringSlice(key string) ([]string, error) {
	return m.valuer.GetStringSlice(key)
}

func (m *Manager) GetDurationSlice(key string) ([]time.Duration, error) {
	return m.valuer.GetDurationSlice(key)
}

// =================== SliceValuer ===================

func (m *Manager) GetStringMap(key string) (map[string]interface{}, error) {
	return m.valuer.GetStringMap(key)
}

func (m *Manager) GetStringMapInt64(key string) (map[string]int64, error) {
	return m.valuer.GetStringMapInt64(key)
}

func (m *Manager) GetStringMapInt(key string) (map[string]int, error) {
	return m.valuer.GetStringMapInt(key)
}

func (m *Manager) GetStringMapBool(key string) (map[string]bool, error) {
	return m.valuer.GetStringMapBool(key)
}

func (m *Manager) GetStringMapString(key string) (map[string]string, error) {
	return m.valuer.GetStringMapString(key)
}

func (m *Manager) GetStringMapStringSlice(key string) (map[string][]string, error) {
	return m.valuer.GetStringMapStringSlice(key)
}
