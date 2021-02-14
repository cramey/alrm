package server

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
)

type worker struct {
	wake  chan bool
	group *config.Group
}

func (w *worker) start(debuglvl int) {
	for {
		if debuglvl > 2 {
			fmt.Printf("%s worker waiting.. \n", w.group.Name)
		}
		<-w.wake
		if debuglvl > 2 {
			fmt.Printf("%s worker wake.. \n", w.group.Name)
		}

		for _, h := range w.group.Hosts {
			for _, c := range h.Checks {
				err := c.Check(debuglvl)
				if err != nil {
					fmt.Printf("check error: %s\n", err)
				}
			}
		}
	}
}
