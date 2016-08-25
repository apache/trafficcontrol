package manager

import (
	"sync"
)

type Event struct {
	Index       uint64 `json:"index"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Hostname    string `json:"hostname"`
	Type        string `json:"type"`
	Available   bool   `json:"isAvailable"`
}

const maxEvents = 200 // TODO make config?

type EventsThreadsafe struct {
	events    *[]Event
	m         *sync.Mutex
	nextIndex *uint64
}

func copyEvents(a []Event) []Event {
	b := make([]Event, len(a), len(a))
	for i, v := range a {
		b[i] = v
	}
	return b
}

func NewEventsThreadsafe() EventsThreadsafe {
	i := uint64(0)
	return EventsThreadsafe{m: &sync.Mutex{}, events: &[]Event{}, nextIndex: &i}
}

func (o *EventsThreadsafe) Get() []Event {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	return copyEvents(*o.events)
}

func (o *EventsThreadsafe) Add(e Event) {
	o.m.Lock()
	e.Index = *o.nextIndex
	*o.nextIndex++
	*o.events = append([]Event{e}, *o.events...)
	if len(*o.events) > maxEvents {
		*o.events = (*o.events)[:maxEvents-1]
	}
	o.m.Unlock()
}
