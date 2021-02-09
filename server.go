package main

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
	"sync"
	"time"
)

func startServer(cfg *config.Config, debuglvl int) error {
	c := sync.NewCond(&sync.Mutex{})
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
		if debuglvl > 2 {
			fmt.Printf("%s goroutine waiting.. \n", g.Name)
		}
		c.L.Lock()
		c.Wait()
		c.L.Unlock()
		if debuglvl > 2 {
			fmt.Printf("%s goroutine wake.. \n", g.Name)
		}

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
