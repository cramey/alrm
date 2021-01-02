package config

import (
	"fmt"
)

type AlrmGroup struct {
	Name  string
	Hosts map[string]*AlrmHost
}

func (ag *AlrmGroup) NewHost(name string) (*AlrmHost, error) {
	if ag.Hosts == nil {
		ag.Hosts = make(map[string]*AlrmHost)
	}

	if _, exists := ag.Hosts[name]; exists {
		return nil, fmt.Errorf("host %s already exists", name)
	}

	host := &AlrmHost{Name: name}
	ag.Hosts[name] = host
	return host, nil
}

func (ag *AlrmGroup) Check(debuglvl int) error {
	for _, host := range ag.Hosts {
		for _, chk := range host.Checks {
			err := chk.Check(debuglvl)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
