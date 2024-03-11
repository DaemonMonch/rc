package rc

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type PollerEvent struct {
	RawEvent interface{}
}

type Poller interface {
	Poll() <-chan PollerEvent
	Close()
}

type FsNotifyPoller struct {
	WatchPath string
	fn        *fsnotify.Watcher
	fnTypes   fsnotify.Op
	c         chan PollerEvent
	close     chan struct{}

	lastEvtTime time.Time
}

func NewFsNotifyPoller(watchPath string, opts fsnotify.Op) *FsNotifyPoller {
	fn, err := fsnotify.NewBufferedWatcher(1)
	if err != nil {
		log.Fatalln(err)
	}
	fn.Add(watchPath)

	return &FsNotifyPoller{WatchPath: watchPath, fn: fn, fnTypes: opts, c: make(chan PollerEvent, 1), close: make(chan struct{})}
}

func (p *FsNotifyPoller) Poll() <-chan PollerEvent {
	go func() {
		for {
			select {
			case e, ok := <-p.fn.Events:
				if ok && e.Has(p.fnTypes) {
					now := time.Now()
					if now.Sub(p.lastEvtTime).Milliseconds() > 10 {
						p.lastEvtTime = now
						p.c <- PollerEvent{RawEvent: e}
					}
				}

			case <-p.close:
				log.Println("FsNotifyPoller close")
				return
			}
		}
	}()
	return p.c
}

func (p *FsNotifyPoller) Close() {
	p.fn.Close()
	close(p.close)
}

type TimedPoller time.Ticker

func (t *TimedPoller) Poll() <-chan PollerEvent {
	timer := ((*time.Ticker)(t))
	c := make(chan PollerEvent, 1)
	go func() {
		for now := range timer.C {
			c <- PollerEvent{RawEvent: now}
		}
		log.Println("TimedPoller close")
	}()

	return c
}

func (t *TimedPoller) Close() {
	((*time.Ticker)(t)).Stop()
}
