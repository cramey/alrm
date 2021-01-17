package config

import (
	"alrm/check"
)

type Host struct {
	Name    string
	Address string
	Checks  []check.Check
}

func (ah *Host) GetAddress() string {
	if ah.Address != "" {
		return ah.Address
	}
	return ah.Name
}

func (ah *Host) NewCheck(name string) (check.Check, error) {
	chk, err := check.NewCheck(name, ah.GetAddress())
	if err != nil {
		return nil, err
	}
	ah.Checks = append(ah.Checks, chk)
	return chk, nil
}
