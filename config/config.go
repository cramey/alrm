package config

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/alarm"
	"os"
	"time"
)

type Config struct {
	Groups     map[string]*Group
	Alarms     map[string]alarm.Alarm
	Interval   time.Duration
	DebugLevel int
	Listen     string
	Path       string
	APIKey     []byte
	APIKeyFile string
}

func (c *Config) NewAlarm(name string, typename string) (alarm.Alarm, error) {
	if c.Alarms == nil {
		c.Alarms = make(map[string]alarm.Alarm)
	}

	if _, exists := c.Alarms[name]; exists {
		return nil, fmt.Errorf("alarm %s already exists", name)
	}

	a, err := alarm.NewAlarm(name, typename)
	if err != nil {
		return nil, err
	}
	c.Alarms[name] = a

	return a, nil
}

func (c *Config) NewGroup(name string) (*Group, error) {
	if c.Groups == nil {
		c.Groups = make(map[string]*Group)
	}

	if _, exists := c.Groups[name]; exists {
		return nil, fmt.Errorf("group %s already exists", name)
	}

	group := &Group{Name: name}
	c.Groups[name] = group
	return group, nil
}

func (c *Config) SetInterval(val string) error {
	interval, err := time.ParseDuration(val)
	if err != nil {
		return err
	}

	c.Interval = interval
	return nil
}

func ReadConfig(fn string, debuglvl int) (*Config, error) {
	cfg := &Config{
		// Default check interval, 30 seconds
		Interval: time.Second * 30,
		// Default listen address
		Listen:     "127.0.0.1:8282",
		DebugLevel: debuglvl,
		Path:       fn,
		// API keyfile defaults to alrmrc.key
		APIKeyFile: fn + ".key",
	}

	pr := &parser{config: cfg}
	if err := pr.parse(); err != nil {
		return nil, err
	}

	if len(cfg.APIKey) == 0 {
		b, err := os.ReadFile(cfg.APIKeyFile)
		if err != nil {
			return nil, err
		}
		cfg.APIKey = b
	}

	return cfg, nil
}
