package server

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
)

type worker struct {
	wakec     chan bool
	shutdownc chan bool
	group     *config.Group
	debuglvl  int
}

func (w *worker) start() {
	for {
		if w.debuglvl > 2 {
			fmt.Printf("%s worker waiting.. \n", w.group.Name)
		}
		<-w.wakec
		if w.debuglvl > 2 {
			fmt.Printf("%s worker wake.. \n", w.group.Name)
		}

		for _, h := range w.group.Hosts {
			for _, c := range h.Checks {
				err := c.Check(w.debuglvl)
				if err != nil {
					fmt.Printf("check error: %s\n", err)
				}
			}
		}
	}
}

// Wake this worker with a non-blocking push
// into the channel
func (w *worker) wake() {
	select {
	case w.wakec <- true:
	default:
	}
}

// Shutdown this worker with a non-blocking push
// into the channel
func (w *worker) shutdown() {
	select {
	case w.shutdownc <- true:
	default:
	}
}

func makeworker(g *config.Group, d int) *worker {
	return &worker{
		group:    g,
		debuglvl: d,
		// This channel is unbuffered so that checks that take
		// over the set interval don't backlog
		wakec: make(chan bool),
		// This channel is buffered because we want it to remember
		// an order to shutdown
		shutdownc: make(chan bool, 1),
	}
}
