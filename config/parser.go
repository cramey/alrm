package config

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/alarm"
	"git.binarythought.com/cdramey/alrm/check"
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

type parser struct {
	config        *Config
	states        []int
	lastHost      *Host
	lastGroup     *Group
	lastCheck     check.Check
	lastAlarm     alarm.Alarm
	lastAlarmName string
}

func (p *parser) parse() error {
	tok, err := NewTokenizer(p.config.Path)
	if err != nil {
		return err
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
				return fmt.Errorf("invalid token in %s, line %d: \"%s\"",
					p.config.Path, tok.Line(), tk)
			}

		case PR_SET:
			key := strings.ToLower(tk)
			if !tok.Scan() {
				return fmt.Errorf("empty value name for set in %s, line %d",
					p.config.Path, tok.Line())
			}

			value := tok.Text()
			switch key {
			case "interval":
				err := p.config.SetInterval(value)
				if err != nil {
					return fmt.Errorf(
						"invalid duration for interval in %s, line %d: \"%s\"",
						p.config.Path, tok.Line(), value,
					)
				}
			case "listen":
				p.config.Listen = value
			case "api.key":
				p.config.APIKey = value
			case "api.keyfile":
				p.config.APIKeyFile = value
			default:
				return fmt.Errorf("unknown key for set in %s, line %d: \"%s\"",
					p.config.Path, tok.Line(), tk,
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
				p.lastGroup, err = p.config.NewGroup(tk)
				if err != nil {
					return fmt.Errorf("%s in %s, line %d",
						err.Error(), p.config.Path, tok.Line(),
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
				p.lastGroup, err = p.config.NewGroup(tk)
				if err != nil {
					return fmt.Errorf("%s in %s, line %d",
						err.Error(), p.config.Path, tok.Line(),
					)
				}
			}

			if p.lastHost == nil {
				p.lastHost, err = p.lastGroup.NewHost(tk)
				if err != nil {
					return fmt.Errorf("%s in %s, line %d",
						err.Error(), p.config.Path, tok.Line(),
					)
				}
				continue
			}

			switch strings.ToLower(tk) {
			case "address":
				if !tok.Scan() {
					return fmt.Errorf("empty address for host in %s, line %d",
						p.config.Path, tok.Line())
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
					return fmt.Errorf("%s in %s, line %d",
						err.Error(), p.config.Path, tok.Line())
				}
				continue
			}
			cont, err := p.lastCheck.Parse(tk)
			if err != nil {
				return fmt.Errorf("%s in %s, line %d",
					err.Error(), p.config.Path, tok.Line())
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

				p.lastAlarm, err = p.config.NewAlarm(p.lastAlarmName, tk)
				if err != nil {
					return fmt.Errorf("%s in %s, line %d",
						err.Error(), p.config.Path, tok.Line())
				}
				p.lastAlarmName = ""
				continue
			}
			cont, err := p.lastAlarm.Parse(tk)
			if err != nil {
				return fmt.Errorf("%s in %s, line %d",
					err.Error(), p.config.Path, tok.Line())
			}
			if !cont {
				p.lastAlarm = nil
				p.prevState()
				goto stateswitch
			}

		default:
			return fmt.Errorf("unknown parser state: %d", p.state())
		}
	}
	if err := tok.Err(); err != nil {
		return err
	}
	return nil
}

func (p *parser) state() int {
	if len(p.states) < 1 {
		return PR_NONE
	}
	return p.states[len(p.states)-1]
}

func (p *parser) setState(state int) {
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

	if p.config.DebugLevel > 1 {
		fmt.Printf("Parser state: %s", p.stateName())
	}
	p.states = append(p.states, state)
	if p.config.DebugLevel > 1 {
		fmt.Printf(" -> %s\n", p.stateName())
	}
}

func (p *parser) prevState() int {
	if len(p.states) > 0 {
		p.states = p.states[:len(p.states)-1]
	}
	return p.state()
}

func (p *parser) stateName() string {
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
