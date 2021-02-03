package alarm

import (
	"fmt"
	"git.binarythought.com/cdramey/alrm/alarm/email"
)

type Alarm interface {
	Parse(string) (bool, error)
	Alarm(string, string, string, error) error
}

func NewAlarm(name string, typename string) (Alarm, error) {
	switch typename {
	case "email":
		return email.NewAlarmEmail(name), nil
	default:
		return nil, fmt.Errorf("unknown alarm name \"%s\"", name)
	}
}
