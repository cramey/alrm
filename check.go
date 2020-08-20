package main

import (
	"fmt"
)

func NewCheck(name string, addr string) (AlrmCheck, error) {
	switch name {
	case "ping":
		return &CheckPing{Address: addr}, nil
	default:
		return nil, fmt.Errorf("unknown check name \"%s\"", name)
	}
}
