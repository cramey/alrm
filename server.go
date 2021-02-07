package main

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
	"time"
	"sync"
)

func startServer(cfg *config.Config, debuglvl int) error {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	for _, g := range cfg.Groups {
		go worker(g, c, debuglvl)
	}

	t := time.NewTicker(cfg.Interval)
	defer t.Stop()
	for {
		select {
		case r := <-t.C:
			if debuglvl > 0 {
				fmt.Printf("Interval check at %s\n", r)
			}
			c.Broadcast()
		}
	}
	return nil
}

func worker(g *config.Group, c *sync.Cond, debuglvl int) {
	for {
		c.L.Lock()
		c.Wait()
		c.L.Unlock()

		for _, h := range g.Hosts {
			for _, c := range h.Checks {
				err := c.Check(debuglvl)
				if err != nil {
					fmt.Printf("Check error: %s\n", err)
				}
			}
		}
	}
}
