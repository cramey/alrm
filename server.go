package main

import (
	"alrm/check"
	"alrm/config"
	"fmt"
	"time"
)

func startServer(cfg *config.Config, debuglvl int) error {
	ch := make(chan check.Check, cfg.Threads)
	for i := 0; i < cfg.Threads; i++ {
		go worker(ch, debuglvl)
	}

	t := time.NewTicker(cfg.Interval)
	defer t.Stop()
	for {
		select {
		case r := <-t.C:
			if debuglvl > 0 {
				fmt.Printf("Interval check at %s\n", r)
			}
			for _, g := range cfg.Groups {
				for _, h := range g.Hosts {
					for _, c := range h.Checks {
						ch <- c
					}
				}
			}
		}
	}
	return nil
}

func worker(ch chan check.Check, debuglvl int) {
	for {
		chk := <-ch
		err := chk.Check(debuglvl)
		if err != nil {
			fmt.Printf("Check error: %s\n", err)
		}
	}
}
