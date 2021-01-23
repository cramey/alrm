package alarm

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"text/template"
	"time"
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
		case "user":
			a.state = TK_USER
		case "pass":
			a.state = TK_PASS
		case "subject":
			a.state = TK_SUBJECT
		case "message":
			a.state = TK_MESSAGE
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
		// If the smtp host doesn't contain a port, add the default
		if !strings.Contains(tk, ":") {
			tk += ":25"
		}
		a.SMTP = tk
		a.state = TK_NONE

	case TK_TO:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("to address cannot be empty")
		}
		a.To = append(a.To, tk)
		a.state = TK_NONE

	case TK_USER:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("user cannot be empty")
		}
		a.User = tk
		a.state = TK_NONE

	case TK_PASS:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("pass cannot be empty")
		}
		a.Pass = tk
		a.state = TK_NONE

	case TK_SUBJECT:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("subject cannot be empty")
		}
		a.Subject = tk
		err := a.updateSubject()
		if err != nil {
			return false, err
		}
		a.state = TK_NONE

	case TK_MESSAGE:
		if strings.TrimSpace(tk) == "" {
			return false, fmt.Errorf("message cannot be empty")
		}
		a.Message = tk
		err := a.updateMessage()
		if err != nil {
			return false, err
		}
		a.state = TK_NONE

	default:
		return false, fmt.Errorf("invalid state in alarm_email")
	}
	return true, nil
}

func (a *AlarmEmail) Alarm(grp string, hst string, chk string, alerr error) error {
	dt := EmailDetail{
		Timestamp: time.Now().Format(time.RFC1123),
		Group:     grp,
		Host:      hst,
		Check:     chk,
		Error:     alerr,
	}

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

	fmt.Fprintf(m, "From: %s\r\n", a.From)
	fmt.Fprintf(m, "To: %s\r\n", strings.Join(a.To, ";"))

	fmt.Fprintf(m, "Subject: ")
	err = a.subTemplate.Execute(m, dt)
	if err != nil {
		return err
	}
	fmt.Fprintf(m, "\r\n\r\n")

	err = a.msgTemplate.Execute(m, dt)
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
