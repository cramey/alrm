package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Command struct {
	Expires   time.Time `json:"exp"`
	Command   string    `json:"cmd"`
	Scheme    string    `json:"sch"`
	Signature []byte    `json:"sig,omitempty"`
}

func ParseCommand(jsn []byte) (*Command, error) {
	cmd := &Command{}
	err := json.Unmarshal(jsn, cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Command) JSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Command) Sign(key []byte) error {
	switch c.Scheme {
	case "hmac-sha256":
		j, err := c.JSON()
		if err != nil {
			return fmt.Errorf("json encoding error")
		}

		mac := hmac.New(sha256.New, key)
		mac.Write(j)
		c.Signature = mac.Sum(nil)

	case "":
		return fmt.Errorf("scheme may not be empty")

	default:
		return fmt.Errorf("unsupported scheme: %s", c.Scheme)
	}
	return nil
}

func (c *Command) Validate(key []byte) error {
	cpy := &Command{
		Expires: c.Expires,
		Command: c.Command,
		Scheme:  c.Scheme,
	}
	err := cpy.Sign(key)
	if err != nil {
		return err
	}

	if !hmac.Equal(cpy.Signature, c.Signature) {
		return fmt.Errorf("invalid signature")
	}

	if time.Now().After(c.Expires) {
		return fmt.Errorf("command expired")
	}

	return nil
}
