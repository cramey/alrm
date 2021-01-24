package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

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
