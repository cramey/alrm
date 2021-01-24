package email

import (
	"fmt"
	"os"
	"text/template"
)

const (
	TK_NONE = iota
	TK_TO
	TK_SMTP
	TK_FROM
	TK_USER
	TK_PASS
	TK_SUBJECT
	TK_MESSAGE
)

type AlarmEmail struct {
	Type        string
	Name        string
	To          []string
	From        string
	SMTP        string
	User        string
	Pass        string
	Subject     string
	subTemplate *template.Template
	Message     string
	msgTemplate *template.Template
	state       int
}

type EmailDetail struct {
	Timestamp string
	Group     string
	Host      string
	Check     string
	Error     error
}

func NewAlarmEmail(name string) *AlarmEmail {
	host, _ := os.Hostname()
	if host == "" {
		host = "localhost"
	}

	al := &AlarmEmail{
		Type:    "email",
		Name:    name,
		From:    "alrm@" + host,
		SMTP:    "localhost:25",
		Subject: "{{.Host}} failure",
		Message: "Check {{.Check}} failed at {{.Timestamp}}: {{.Error}}",
	}
	al.updateSubject()
	al.updateMessage()
	return al
}

func (a *AlarmEmail) updateSubject() error {
	t := template.New("email subject")
	_, err := t.Parse(a.Subject)
	if err != nil {
		fmt.Print(err)
		return err
	}
	a.subTemplate = t
	return nil
}

func (a *AlarmEmail) updateMessage() error {
	t := template.New("email message")
	_, err := t.Parse(a.Message)
	if err != nil {
		fmt.Print(err)
		return err
	}
	a.msgTemplate = t
	return nil
}
