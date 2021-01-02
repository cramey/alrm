package check

import (
	"fmt"
)

type CheckPing struct {
	Type string
	Address string
}

func (c *CheckPing) Check() error {
	fmt.Printf("Pinging %s .. \n", c.Address)
	return nil
}

func (c *CheckPing) Parse(tk string) (bool, error) {
	return false, nil
}
