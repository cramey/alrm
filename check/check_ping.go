package check

type CheckPing struct {
	Type string
	Address string
}

func (c *CheckPing) Check() error {
	return nil
}

func (c *CheckPing) Parse(tk string) (bool, error) {
	return false, nil
}
