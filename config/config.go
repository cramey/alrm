package config

import (
	"fmt"
	"strconv"
)

type AlrmConfig struct {
	Groups   map[string]*AlrmGroup
	Interval int
}

func (ac *AlrmConfig) NewGroup(name string) (*AlrmGroup, error) {
	if ac.Groups == nil {
		ac.Groups = make(map[string]*AlrmGroup)
	}

	if _, exists := ac.Groups[name]; exists {
		return nil, fmt.Errorf("group %s already exists", name)
	}

	group := &AlrmGroup{Name: name}
	ac.Groups[name] = group
	return group, nil
}

func (ac *AlrmConfig) SetInterval(val string) error {
	interval, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	ac.Interval = interval
	return nil
}

func ReadConfig(fn string, debuglvl int) (*AlrmConfig, error) {
	parser := &Parser{DebugLevel: debuglvl}
	config, err := parser.Parse(fn)
	if err != nil {
		return nil, err
	}
	return config, nil
}
