package docker

import (
	eventtypes "github.com/docker/docker/api/types/events"
	"sync"
)

type eventHandler struct {
	handlers map[string]func(eventtypes.Message)
	mu       sync.Mutex
}

func (w *eventHandler) Handle(action string, h func(eventtypes.Message)) {
	w.mu.Lock()
	w.handlers[action] = h
	w.mu.Unlock()
}

// Watch ranges over the passed in event chan and processes the events based on the
// handlers created for a given action.
// To stop watching, close the event chan.
func (w *eventHandler) Watch(c <-chan eventtypes.Message) {
	for e := range c {
		w.mu.Lock()
		h, existsH := w.handlers[e.Action]
		i, existsI := w.handlers["*"]
		w.mu.Unlock()
		if existsH {
			go h(e)
		}
		if existsI {
			go i(e)
		}

	}
}
