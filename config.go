package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type AlrmConfig struct {
	Groups []*AlrmGroup
}

func (ac *AlrmConfig) NewGroup() *AlrmGroup {
	group := &AlrmGroup{}
	ac.Groups = append(ac.Groups, group)
	return group
}

func (ac *AlrmConfig) LastGroup() *AlrmGroup {
	return ac.Groups[len(ac.Groups)-1]
}

type AlrmGroup struct {
	Name  string
	Hosts []*AlrmHost
}

func (ag *AlrmGroup) NewHost() *AlrmHost {
	host := &AlrmHost{}
	ag.Hosts = append(ag.Hosts, host)
	return host
}

func (ag *AlrmGroup) LastHost() *AlrmHost {
	return ag.Hosts[len(ag.Hosts)-1]
}

type AlrmHost struct {
	Name    string
	Address string
}

const (
	TK_NONE = iota
	TK_MONITOR
	TK_GROUP
	TK_HOST
)

func stateName(s int) string {
	switch s {
	case TK_NONE:
		return "TK_NONE"
	case TK_MONITOR:
		return "TK_MONTIOR"
	case TK_GROUP:
		return "TK_GROUP"
	case TK_HOST:
		return "TK_HOST"
	default:
		return "UNKNOWN"
	}
}

func ReadConfig(fn string) (*AlrmConfig, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("cannot open config \"%s\": %s", fn, err.Error())
	}
	defer file.Close()

	config := &AlrmConfig{}
	parser := &Parser{}

	scan := bufio.NewScanner(file)
	scan.Split(parser.Split)
	for scan.Scan() {
		tk := scan.Text()

	stateswitch:
		switch parser.GetState() {
		case TK_NONE:
			switch strings.ToLower(tk) {
			case "monitor":
				parser.SetState(TK_MONITOR)
			default:
				return nil, fmt.Errorf("Invalid token in %s, line %d: \"%s\"",
					fn, parser.Line+1, tk)
			}

		case TK_MONITOR:
			switch strings.ToLower(tk) {
			case "host":
				config.NewGroup().NewHost()
				parser.SetState(TK_HOST)

			case "group":
				config.NewGroup()
				parser.SetState(TK_GROUP)

			default:
				parser.PrevState()
				goto stateswitch
			}

		case TK_GROUP:
			group := config.LastGroup()

			switch strings.ToLower(tk) {
			case "host":
				group.NewHost()
				parser.SetState(TK_HOST)
				continue

			default:
				if group.Name == "" {
					group.Name = tk
					continue
				}

				parser.PrevState()
				goto stateswitch
			}

		case TK_HOST:
			host := config.LastGroup().LastHost()
			if host.Name == "" {
				host.Name = tk
				continue
			}

			switch strings.ToLower(tk) {
			case "address":
				if scan.Scan() {
					host.Address = scan.Text()
				}
				continue

			default:
				parser.PrevState()
				goto stateswitch
			}

		default:
			return nil, fmt.Errorf("Unknown parser state: %d", parser.GetState())
		}
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return config, nil
}

type Parser struct {
	Line  int
	states []int
}

func (pr *Parser) GetState() int {
	if len(pr.states) < 1 {
		return TK_NONE
	}
	return pr.states[len(pr.states) - 1]
}

func (pr *Parser) SetState(state int) {
	//fmt.Printf("%s -> %s\n", stateName(pr.GetState()), stateName(state))
	pr.states = append(pr.states, state)
}

func (pr *Parser) PrevState() int {
	if len(pr.states) > 0 {
		pr.states = pr.states[:len(pr.states) - 1]
	}
	return pr.GetState()
}

func (pr *Parser) Split(data []byte, atEOF bool) (int, []byte, error) {
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
			pr.Line++
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
