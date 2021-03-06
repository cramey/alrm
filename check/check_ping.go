package check

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/check/ping"
	"strconv"
	"strings"
	"time"
)

const (
	TK_NONE = iota
	TK_COUNT
	TK_TIMEOUT
	TK_INTERVAL
	TK_MAXPACKETLOSS
	TK_MAXRTT
	TK_MAXAVGRTT
)

type CheckPing struct {
	Type          string
	Address       string
	Count         int
	Timeout       time.Duration
	Interval      time.Duration
	MaxPacketLoss float64
	MaxRTT        time.Duration
	MaxAvgRTT     time.Duration
	state         int
}

func NewCheckPing(addr string) *CheckPing {
	return &CheckPing{
		Type: "ping", Address: addr,
		Count: 1, Timeout: time.Second * 5,
		Interval: time.Second, MaxPacketLoss: 0,
	}
}

func (c *CheckPing) Check(debuglvl int) error {
	if debuglvl > 0 {
		fmt.Printf("Pinging %s .. \n", c.Address)
	}

	p, err := ping.NewPinger(c.Address)
	if err != nil {
		return err
	}

	p.Count = c.Count
	p.Timeout = c.Timeout
	p.Interval = c.Interval

	err = p.Run()
	if err != nil {
		return err
	}

	stats := p.Statistics()

	if debuglvl > 1 {
		for i, r := range stats.Rtts {
			fmt.Printf("Ping %d: %s\n", i+1, r)
		}
	}

	if debuglvl > 1 {
		fmt.Printf("Packet loss: %.f%%\n", stats.PacketLoss)
		fmt.Printf("Max RTT: %s\n", stats.MaxRtt.String())
		fmt.Printf("Max Configured RTT: %s\n", c.MaxRTT.String())
		fmt.Printf("Average RTT: %s\n", stats.AvgRtt.String())
	}

	if stats.PacketLoss > c.MaxPacketLoss {
		return fmt.Errorf("ping packet loss exceeds max")
	}

	if c.MaxRTT > 0 && stats.MaxRtt > c.MaxRTT {
		return fmt.Errorf("ping rtt exceeds max")
	}

	if c.MaxAvgRTT > 0 && stats.AvgRtt > c.MaxAvgRTT {
		return fmt.Errorf("ping average rtt exceeds max")
	}
	return nil
}

func (c *CheckPing) Parse(tk string) (bool, error) {
	var err error
	switch c.state {
	case TK_NONE:
		switch strings.ToLower(tk) {
		case "count":
			c.state = TK_COUNT
		case "timeout":
			c.state = TK_TIMEOUT
		case "interval":
			c.state = TK_INTERVAL
		case "maxpacketloss":
			c.state = TK_MAXPACKETLOSS
		case "maxrtt":
			c.state = TK_MAXRTT
		case "maxavgrtt":
			c.state = TK_MAXAVGRTT
		default:
			return false, nil
		}

	case TK_COUNT:
		c.Count, err = strconv.Atoi(tk)
		if err != nil {
			return false, fmt.Errorf("invalid count \"%s\"", tk)
		}
		c.state = TK_NONE

	case TK_TIMEOUT:
		c.Timeout, err = time.ParseDuration(tk)
		if err != nil {
			return false, fmt.Errorf("invalid timeout \"%s\"", tk)
		}
		c.state = TK_NONE

	case TK_INTERVAL:
		c.Interval, err = time.ParseDuration(tk)
		if err != nil {
			return false, fmt.Errorf("invalid interval \"%s\"", tk)
		}
		c.state = TK_NONE

	case TK_MAXPACKETLOSS:
		c.MaxPacketLoss, err = strconv.ParseFloat(tk, 64)
		if err != nil {
			return false, fmt.Errorf("invalid minpacketloss \"%s\"", tk)
		}
		c.state = TK_NONE

	case TK_MAXRTT:
		c.MaxRTT, err = time.ParseDuration(tk)
		if err != nil {
			return false, fmt.Errorf("invalid maxrtt \"%s\"", tk)
		}
		c.state = TK_NONE

	case TK_MAXAVGRTT:
		c.MaxAvgRTT, err = time.ParseDuration(tk)
		if err != nil {
			return false, fmt.Errorf("invalid maxavgrtt \"%s\"", tk)
		}
		c.state = TK_NONE

	default:
		return false, fmt.Errorf("invalid state in check_ping")
	}
	return true, nil
}
