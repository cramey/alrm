package check

import (
	"fmt"
)

type Check interface {
	Parse(string) (bool, error)
	Check(int) error
}

func NewCheck(name string, addr string) (Check, error) {
	switch name {
	case "ping":
		return NewCheckPing(addr), nil
	default:
		return nil, fmt.Errorf("unknown check name \"%s\"", name)
	}
}
