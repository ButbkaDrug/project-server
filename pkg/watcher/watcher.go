package watcher

import (
	"os"
	"time"
)

type Watcher struct {
	list   map[string]time.Time
	Reload *chan bool
}

// returnes new instance of a Watcher.
func New() *Watcher {
	ch := make(chan bool)
	return &Watcher{
		list:   make(map[string]time.Time),
		Reload: &ch,
	}
}

// adds an entity to a list
func (w *Watcher) Add(id string, t time.Time) {
	w.list[id] = t
}

// will run an infinite loop checking for changes in the files
// if change found will send true to a Reload chan
func (w *Watcher) Watch() {
	//FIX: get rid of infinate loop
	for {
		w.listCheck()
		time.Sleep(time.Second)
	}
}

func (w *Watcher) listCheck() {
	for key, val := range w.list {
		s, err := os.Stat(key)
		if err != nil {
			//TODO: log the error and move on.
			continue
		}

		if s.ModTime() != val {
			w.list[key] = s.ModTime()
			if !s.ModTime().IsZero() {
				*w.Reload <- true
			}
		}
	}
}
