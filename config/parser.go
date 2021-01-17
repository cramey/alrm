package config

import (
	"alrm/check"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	TK_NONE = iota
	TK_SET
	TK_MONITOR
	TK_GROUP
	TK_HOST
	TK_CHECK
)

type Parser struct {
	DebugLevel int
	Line       int
	states     []int
	lasthost   *AlrmHost
	lastgroup  *AlrmGroup
	lastcheck  check.AlrmCheck
}

func (p *Parser) Parse(fn string) (*AlrmConfig, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("cannot open config \"%s\": %s", fn, err.Error())
	}
	defer file.Close()

	config := NewConfig()

	scan := bufio.NewScanner(file)
	scan.Split(p.Split)
	for scan.Scan() {
		tk := scan.Text()
	stateswitch:
		switch p.state() {
		case TK_NONE:
			switch strings.ToLower(tk) {
			case "monitor":
				p.setState(TK_MONITOR)
			case "set":
				p.setState(TK_SET)
			default:
				return nil, fmt.Errorf("invalid token in %s, line %d: \"%s\"",
					fn, p.Line, tk)
			}

		case TK_SET:
			key := strings.ToLower(tk)
			if !scan.Scan() {
				return nil, fmt.Errorf("empty value name for set in %s, line %d",
					fn, p.Line)
			}

			value := scan.Text()
			switch key {
			case "interval":
				err := config.SetInterval(value)
				if err != nil {
					return nil, fmt.Errorf(
						"invalid number for interval in %s, line %d: \"%s\"",
						fn, p.Line, value,
					)
				}
			default:
				return nil, fmt.Errorf("unknown key for set in %s, line %d: \"%s\"",
					fn, p.Line, tk,
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
			if p.lastgroup == nil {
				p.lastgroup, err = config.NewGroup(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, p.Line,
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
			if p.lastgroup == nil {
				p.lastgroup, err = config.NewGroup(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, p.Line,
					)
				}
			}

			if p.lasthost == nil {
				p.lasthost, err = p.lastgroup.NewHost(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, p.Line,
					)
				}
				continue
			}

			switch strings.ToLower(tk) {
			case "address":
				if !scan.Scan() {
					return nil, fmt.Errorf("empty address for host in %s, line %d",
						fn, p.Line)
				}
				p.lasthost.Address = scan.Text()

			case "check":
				p.setState(TK_CHECK)

			default:
				p.prevState()
				goto stateswitch
			}

		case TK_CHECK:
			if p.lastcheck == nil {
				p.lastcheck, err = p.lasthost.NewCheck(tk)
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, p.Line)
				}
				continue
			}
			cont, err := p.lastcheck.Parse(tk)
			if err != nil {
				return nil, fmt.Errorf("%s in %s, line %d",
					err.Error(), fn, p.Line)
			}
			if !cont {
				p.lastcheck = nil
				p.prevState()
				goto stateswitch
			}

		default:
			return nil, fmt.Errorf("unknown parser state: %d", p.state())
		}
	}
	if err := scan.Err(); err != nil {
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
		p.lastgroup = nil
		fallthrough
	case TK_HOST:
		p.lasthost = nil
		p.lastcheck = nil
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
		return "TK_MONTIOR"
	case TK_GROUP:
		return "TK_GROUP"
	case TK_HOST:
		return "TK_HOST"
	case TK_CHECK:
		return "TK_CHECK"
	default:
		return "UNKNOWN"
	}
}

func (p *Parser) Split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var ignoreline bool
	var started bool
	var startidx int
	var quote byte

	for i := 0; i < len(data); i++ {
		c := data[i]
		// fmt.Printf("%c (%t) (%t)\n", c, started, ignoreline)
		switch c {
		case '\f', '\n', '\r':
			p.Line++
			if ignoreline {
				ignoreline = false
				continue
			}
			fallthrough

		case ' ', '\t', '\v':
			if started && quote == 0 {
				return i + 1, data[startidx:i], nil
			}

		case '\'', '"', '`':
			if started && quote == c {
				return i + 1, data[startidx:i], nil
			}

			if !ignoreline && quote == 0 {
				quote = c
			}

		case '#':
			if !started {
				ignoreline = true
			}

		default:
			if !ignoreline && !started {
				started = true
				startidx = i
			}
		}
	}

	if atEOF {
		if ignoreline {
			return len(data), nil, nil
		}
		if started {
			return len(data), data[startidx:], nil
		}
	}

	return 0, nil, nil
}
