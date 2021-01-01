package config

import (
	"alrm/check"
)

type AlrmHost struct {
	Name    string
	Address string
	Checks  []check.AlrmCheck
}

func (ah *AlrmHost) GetAddress() string {
	if ah.Address != "" {
		return ah.Address
	}
	return ah.Name
}

func (ah *AlrmHost) NewCheck(name string) (check.AlrmCheck, error) {
	chk, err := check.NewCheck(name, ah.GetAddress())
	if err != nil {
		return nil, err
	}
	ah.Checks = append(ah.Checks, chk)
	return chk, nil
}
