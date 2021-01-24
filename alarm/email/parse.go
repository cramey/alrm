package email

import (
	"fmt"
	"strings"
)

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
