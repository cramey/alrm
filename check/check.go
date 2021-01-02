package check

import (
	"fmt"
)

type AlrmCheck interface {
	Parse(string) (bool, error)
	Check(int) error
}

func NewCheck(name string, addr string) (AlrmCheck, error) {
	switch name {
	case "ping":
		return NewCheckPing(addr), nil
	default:
		return nil, fmt.Errorf("unknown check name \"%s\"", name)
	}
}
