package config

import (
	"fmt"
)

type Group struct {
	Name  string
	Hosts map[string]*Host
}

func (ag *Group) NewHost(name string) (*Host, error) {
	if ag.Hosts == nil {
		ag.Hosts = make(map[string]*Host)
	}

	if _, exists := ag.Hosts[name]; exists {
		return nil, fmt.Errorf("host %s already exists", name)
	}

	host := &Host{Name: name}
	ag.Hosts[name] = host
	return host, nil
}

func (ag *Group) Check(debuglvl int) error {
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
