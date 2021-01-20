package alarm

import (
	"fmt"
	"net/smtp"
	"os"
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
	host, _ := os.Hostname()
	if host == "" {
		host = "localhost"
	}

	return &AlarmEmail{
		Type: "email", Name: name,
		From: "alrm@" + host, SMTP: "localhost:25",
	}
}

func (a *AlarmEmail) Alarm() error {
	c, err := smtp.Dial(a.SMTP)
	if err != nil {
		return err
	}

	helo := "localhost"
	tspl := strings.Split(a.From, "@")
	if len(tspl) > 1 {
		helo = tspl[1]
	}

	err = c.Hello(helo)
	if err != nil {
		return err
	}
	err = c.Mail(a.From)
	if err != nil {
		return err
	}
	for _, to := range a.To {
		err = c.Rcpt(to)
		if err != nil {
			return err
		}
	}
	m, err := c.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\r\n", a.From)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(a.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", "test subject")

	_, err = fmt.Fprintf(m, "%s", msg)
	if err != nil {
		return err
	}
	err = m.Close()
	if err != nil {
		return err
	}
	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}

func (a *AlarmEmail) Parse(tk string) (bool, error) {
	switch a.state {
	case TK_NONE:
		switch strings.ToLower(tk) {
		case "to":
			a.state = TK_TO
		case "from":
			a.state = TK_FROM
		case "smtp":
			a.state = TK_SMTP
		default:
			if len(a.To) < 1 {
				return false, fmt.Errorf("email alarm requires to address")
			}
			return false, nil
		}

	case TK_FROM:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("from address cannot be empty")
		}
		a.From = tk
		a.state = TK_NONE

	case TK_SMTP:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("smtp server cannot be empty")
		}
		a.SMTP = tk
		a.state = TK_NONE

	case TK_TO:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("to address cannot be empty")
		}
		a.To = append(a.To, tk)
		a.state = TK_NONE

	default:
		return false, fmt.Errorf("invalid state in alarm_email")
	}
	return true, nil
}
