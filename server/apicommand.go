package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type APICommand struct {
	Expires   time.Time `json:"exp"`
	Command   string    `json:"cmd"`
	Scheme    string    `json:"sch"`
	Signature []byte    `json:"sig,omitempty"`
}

func ParseAPICommand(jsn []byte) (*APICommand, error) {
	api := &APICommand{}
	err := json.Unmarshal(jsn, api)
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (ac *APICommand) JSON() ([]byte, error) {
	return json.Marshal(ac)
}

func (ac *APICommand) Sign(key []byte) error {
	switch ac.Scheme {
	case "hmac-sha256":
		j, err := ac.JSON()
		if err != nil {
			return fmt.Errorf("json encoding error")
		}

		mac := hmac.New(sha256.New, key)
		mac.Write(j)
		ac.Signature = mac.Sum(nil)

	case "":
		return fmt.Errorf("scheme may not be empty")

	default:
		return fmt.Errorf("unsupported scheme: %s", ac.Scheme)
	}
	return nil
}

func (ac *APICommand) Validate(key []byte) error {
	cpy := &APICommand{
		Expires: ac.Expires,
		Command: ac.Command,
		Scheme:  ac.Scheme,
	}
	err := cpy.Sign(key)
	if err != nil {
		return err
	}

	if !hmac.Equal(cpy.Signature, ac.Signature) {
		return fmt.Errorf("invalid signature")
	}

	if time.Now().After(ac.Expires) {
		return fmt.Errorf("command expired")
	}

	return nil
}
