package easyconfmgr

// Watcher  monitors whether the data source has changed and, if so, notifies the changed event
type Watcher interface {
	// Watch start watch
	Watch() error
	// Events event sent to channel
	Events() <-chan *Event
	// Stop stop watch
	Stop() error
}
type EventKey interface {
	KeyName() string
}

type Event struct {
	description interface{ String() string }
	data        []byte
}

func (e *Event) Data() []byte {
	return e.data
}

func (e *Event) String() string {
	return e.description.String() + "/n" + string(e.data)
}

func NewEvent(description interface{ String() string }, data []byte) *Event {
	return &Event{description: description, data: data}
}
