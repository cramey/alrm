package config

import (
	"alrm/alarm"
	"alrm/check"
	"fmt"
	"strings"
)

const (
	PR_NONE = iota
	PR_SET
	PR_MONITOR
	PR_GROUP
	PR_HOST
	PR_CHECK
	PR_ALARM
)

type Parser struct {
	DebugLevel    int
	states        []int
	lastHost      *Host
	lastGroup     *Group
	lastCheck     check.Check
	lastAlarm     alarm.Alarm
	lastAlarmName string
}

func (p *Parser) Parse(fn string) (*Config, error) {
	config := NewConfig()
	tok, err := NewTokenizer(fn)
	if err != nil {
		return nil, err
	}
	defer tok.Close()

	for tok.Scan() {
		tk := tok.Text()
	stateswitch:
		switch p.state() {
		case PR_NONE:
			switch strings.ToLower(tk) {
			case "monitor":
				p.setState(PR_MONITOR)
			case "set":
				p.setState(PR_SET)
			case "alarm":
				p.setState(PR_ALARM)
			default:
				return nil, fmt.Errorf("invalid token in %s, line %d: \"%s\"",
					fn, tok.Line(), tk)
			}

		case PR_SET:
			key := strings.ToLower(tk)
			if !tok.Scan() {
				return nil, fmt.Errorf("empty value name for set in %s, line %d",
					fn, tok.Line())
			}

			value := tok.Text()
			switch key {
			case "interval":
				err := config.SetInterval(value)
				if err != nil {
					return nil, fmt.Errorf(
						"invalid duration for interval in %s, line %d: \"%s\"",
						fn, tok.Line(), value,
					)
				}
			case "threads":
				err := config.SetThreads(value)
				if err != nil {
					return nil, fmt.Errorf(
						"invalid number for interval in %s, line %d: \"%s\"",
						fn, tok.Line(), value,
					)
				}
			default:
				return nil, fmt.Errorf("unknown key for set in %s, line %d: \"%s\"",
					fn, tok.Line(), tk,
				)
			}
			p.prevState()

		case PR_MONITOR:
			switch strings.ToLower(tk) {
			case "host":
				p.setState(PR_HOST)

			case "group":
				p.setState(PR_GROUP)

			default:
				p.prevState()
				goto stateswitch
			}

		case PR_GROUP:
			if p.lastGroup == nil {
				p.lastGroup, err = config.NewGroup(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, tok.Line(),
					)
				}
				continue
			}

			switch strings.ToLower(tk) {
			case "host":
				p.setState(PR_HOST)

			default:
				p.prevState()
				goto stateswitch
			}

		case PR_HOST:
			// If a host has no group, inherit the host name
			if p.lastGroup == nil {
				p.lastGroup, err = config.NewGroup(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, tok.Line(),
					)
				}
			}

			if p.lastHost == nil {
				p.lastHost, err = p.lastGroup.NewHost(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, tok.Line(),
					)
				}
				continue
			}

			switch strings.ToLower(tk) {
			case "address":
				if !tok.Scan() {
					return nil, fmt.Errorf("empty address for host in %s, line %d",
						fn, tok.Line())
				}
				p.lastHost.Address = tok.Text()

			case "check":
				p.setState(PR_CHECK)

			default:
				p.prevState()
				goto stateswitch
			}

		case PR_CHECK:
			if p.lastCheck == nil {
				p.lastCheck, err = p.lastHost.NewCheck(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, tok.Line())
				}
				continue
			}
			cont, err := p.lastCheck.Parse(tk)
			if err != nil {
				return nil, fmt.Errorf("%s in %s, line %d",
					err.Error(), fn, tok.Line())
			}
			if !cont {
				p.lastCheck = nil
				p.prevState()
				goto stateswitch
			}

		case PR_ALARM:
			if p.lastAlarm == nil {
				if p.lastAlarmName == "" {
					p.lastAlarmName = tk
					continue
				}

				p.lastAlarm, err = config.NewAlarm(p.lastAlarmName, tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, tok.Line())
				}
				p.lastAlarmName = ""
				continue
			}
			cont, err := p.lastAlarm.Parse(tk)
			if err != nil {
				return nil, fmt.Errorf("%s in %s, line %d",
					err.Error(), fn, tok.Line())
			}
			if !cont {
				p.lastAlarm = nil
				p.prevState()
				goto stateswitch
			}

		default:
			return nil, fmt.Errorf("unknown parser state: %d", p.state())
		}
	}
	if err := tok.Err(); err != nil {
		return nil, err
	}
	return config, nil
}

func (p *Parser) state() int {
	if len(p.states) < 1 {
		return PR_NONE
	}
	return p.states[len(p.states)-1]
}

func (p *Parser) setState(state int) {
	switch state {
	case PR_SET, PR_MONITOR:
		fallthrough
	case PR_GROUP:
		p.lastGroup = nil
		fallthrough
	case PR_HOST:
		p.lastHost = nil
		p.lastCheck = nil
	}

	if p.DebugLevel > 1 {
		fmt.Printf("Parser state: %s", p.stateName())
	}
	p.states = append(p.states, state)
	if p.DebugLevel > 1 {
		fmt.Printf(" -> %s\n", p.stateName())
	}
}

func (p *Parser) prevState() int {
	if len(p.states) > 0 {
		p.states = p.states[:len(p.states)-1]
	}
	return p.state()
}

func (p *Parser) stateName() string {
	switch p.state() {
	case PR_NONE:
		return "PR_NONE"
	case PR_SET:
		return "PR_SET"
	case PR_MONITOR:
		return "PR_MONITOR"
	case PR_GROUP:
		return "PR_GROUP"
	case PR_HOST:
		return "PR_HOST"
	case PR_CHECK:
		return "PR_CHECK"
	case PR_ALARM:
		return "PR_ALARM"
	default:
		return "UNKNOWN"
	}
}
