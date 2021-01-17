package alarm

import (
	"fmt"
	"strings"
)

const (
	TK_NONE = iota
  TK_TO
  TK_SMTP
  TK_FROM
)


type AlarmEmail struct {
	Type  string
	Name  string
	From  string
	SMTP  string
	To    []string
	state int
}

func NewAlarmEmail(name string) *AlarmEmail {
	return &AlarmEmail{
		Type: "email", Name: name,
	}
}

func (a *AlarmEmail) Alarm() error {
	fmt.Printf("email alarm")
	return nil
}

func (a *AlarmEmail) Parse(tk string) (bool, error) {
	switch a.state {
		case TK_NONE:
			switch strings.ToLower(tk){
				case "to":
					a.state = TK_TO
				case "from":
					a.state = TK_FROM
				case "smtp":
					a.state = TK_SMTP
				default:
					return false, nil
			}

		case TK_FROM:
			a.From = tk
			a.state = TK_NONE

		case TK_SMTP:
			a.SMTP = tk
			a.state = TK_NONE

		case TK_TO:
			a.To = append(a.To, tk)
			a.state = TK_NONE

		default:
			return false, fmt.Errorf("invalid state in alarm_email")
	}
	return true, nil
}
