package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	Line   int
	states []int
}

func (p *Parser) Parse(fn string) (*AlrmConfig, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("cannot open config \"%s\": %s", fn, err.Error())
	}
	defer file.Close()

	config := &AlrmConfig{}
	var group *AlrmGroup
	var host *AlrmHost
	var check AlrmCheck

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
					fn, p.Line+1, tk)
			}

		case TK_SET:
			key := strings.ToLower(tk)
			if !scan.Scan() {
				return nil, fmt.Errorf("empty value name for set in %s, line %d",
					fn, p.Line+1)
			}

			value := scan.Text()
			switch key {
			case "interval":
				config.Interval, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf(
						"invalid number for interval in %s, line %d: \"%s\"",
						fn, p.Line+1, value,
					)
				}
			default:
				return nil, fmt.Errorf("unknown key for set in %s, line %d: \"%s\"",
					fn, p.Line+1, tk,
				)
			}
			p.prevState()

		case TK_MONITOR:
			switch strings.ToLower(tk) {
			case "host":
				host = config.NewGroup().NewHost()
				p.setState(TK_HOST)

			case "group":
				group = config.NewGroup()
				p.setState(TK_GROUP)

			default:
				p.prevState()
				goto stateswitch
			}

		case TK_GROUP:
			if group == nil {
				return nil, fmt.Errorf("group without initialization")
			}

			switch strings.ToLower(tk) {
			case "host":
				host = group.NewHost()
				p.setState(TK_HOST)
				continue

			default:
				if group.Name == "" {
					group.Name = tk
					continue
				}

				p.prevState()
				goto stateswitch
			}

		case TK_HOST:
			if host == nil {
				return nil, fmt.Errorf("host token without initialization")
			}

			if host.Name == "" {
				host.Name = tk
				continue
			}

			switch strings.ToLower(tk) {
			case "address":
				if !scan.Scan() {
					return nil, fmt.Errorf("empty address for host in %s, line %d",
						fn, p.Line+1)
				}
				host.Address = scan.Text()

			case "check":
				check = nil
				p.setState(TK_CHECK)

			default:
				p.prevState()
				goto stateswitch
			}

		case TK_CHECK:
			if check == nil {
				if host == nil {
					return nil, fmt.Errorf("host token without initialization")
				}
				check, err = NewCheck(strings.ToLower(tk), host.GetAddress())
				if err != nil {
					return nil, fmt.Errorf("%s in %s, line %d",
						err.Error(), fn, p.Line+1)
				}
				host.Checks = append(host.Checks, check)
				continue
			}
			check.Parse(tk)

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
	// fmt.Printf("%s", p.stateName())
	p.states = append(p.states, state)
	// fmt.Printf(" -> %s\n", p.stateName())
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

			if quote == 0 {
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
