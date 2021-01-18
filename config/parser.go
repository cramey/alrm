package config

import (
	"alrm/alarm"
	"alrm/check"
	"fmt"
	"strings"
)

const (
	TK_NONE = iota
	TK_SET
	TK_MONITOR
	TK_GROUP
	TK_HOST
	TK_CHECK
	TK_ALARM
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
		case TK_NONE:
			switch strings.ToLower(tk) {
			case "monitor":
				p.setState(TK_MONITOR)
			case "set":
				p.setState(TK_SET)
			case "alarm":
				p.setState(TK_ALARM)
			default:
				return nil, fmt.Errorf("invalid token in %s, line %d: \"%s\"",
					fn, tok.Line(), tk)
			}

		case TK_SET:
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

		case TK_MONITOR:
			switch strings.ToLower(tk) {
			case "host":
				p.setState(TK_HOST)

			case "group":
				p.setState(TK_GROUP)

			default:
				p.prevState()
				goto stateswitch
			}

		case TK_GROUP:
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
				p.setState(TK_HOST)

			default:
				p.prevState()
				goto stateswitch
			}

		case TK_HOST:
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
				p.setState(TK_CHECK)

			default:
				p.prevState()
				goto stateswitch
			}

		case TK_CHECK:
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

		case TK_ALARM:
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
		return TK_NONE
	}
	return p.states[len(p.states)-1]
}

func (p *Parser) setState(state int) {
	switch state {
	case TK_SET, TK_MONITOR:
		fallthrough
	case TK_GROUP:
		p.lastGroup = nil
		fallthrough
	case TK_HOST:
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
	case TK_NONE:
		return "TK_NONE"
	case TK_SET:
		return "TK_SET"
	case TK_MONITOR:
		return "TK_MONITOR"
	case TK_GROUP:
		return "TK_GROUP"
	case TK_HOST:
		return "TK_HOST"
	case TK_CHECK:
		return "TK_CHECK"
	case TK_ALARM:
		return "TK_ALARM"
	default:
		return "UNKNOWN"
	}
}
