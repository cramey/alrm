package check

import (
	"alrm/check/ping"
	"fmt"
	"time"
)

type CheckPing struct {
	Type       string
	Address    string
	Count      int
	Timeout    time.Duration
}

func (c *CheckPing) Check(debuglvl int) error {
	if debuglvl > 0 {
		fmt.Printf("Pinging %s .. \n", c.Address)
	}

	p, err := ping.NewPinger(c.Address)
	if err != nil {
		return err
	}

	p.Count = 1
	p.Timeout = time.Second * 5
	err = p.Run()
	if err != nil {
		return err
	}

	stats := p.Statistics()
	if len(stats.Rtts) < 1 {
		return fmt.Errorf("ping failure")
	}

	if debuglvl > 0 {
		fmt.Printf("Ping RTT: %s\n", stats.Rtts[0])
	}
	return nil
}

func (c *CheckPing) Parse(tk string) (bool, error) {
	return false, nil
}
