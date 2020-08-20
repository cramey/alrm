package main

type AlrmConfig struct {
	Groups   []*AlrmGroup
	Interval int
}

func (ac *AlrmConfig) NewGroup() *AlrmGroup {
	group := &AlrmGroup{}
	ac.Groups = append(ac.Groups, group)
	return group
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

type AlrmHost struct {
	Name    string
	Address string
	Checks  []AlrmCheck
}

func (ah *AlrmHost) GetAddress() string {
	if ah.Address != "" {
		return ah.Address
	}
	return ah.Name
}

type AlrmCheck interface {
	Parse(string) (bool, error)
	Check() error
}

func ReadConfig(fn string) (*AlrmConfig, error) {
	parser := &Parser{}
	config, err := parser.Parse(fn)
	if err != nil {
		return nil, err
	}
	return config, nil
}
