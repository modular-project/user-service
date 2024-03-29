package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"sync"
	"text/template"
	"users-service/pkg"
)

var (
	app string
)

type jsonEmail struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type send struct {
	App  string
	Code string
}

type Mail struct {
}

var (
	auth     smtp.Auth
	config   *jsonEmail
	onceMail sync.Once
)

func NewMail() Mail {
	return Mail{}
}

// LoadFiles .
func LoadMail() error {
	var err error
	onceMail.Do(func() {
		err = loadData()
	})

	return err
}

func getEnv(env string) (string, error) {
	s, f := os.LookupEnv(env)
	if !f {
		return "", fmt.Errorf("environment variable (%s) not found", env)
	}
	return s, nil
}

func loadData() error {
	user, err := getEnv("MAIL_USER")
	if err != nil {
		return err
	}
	password, err := getEnv("MAIL_PWD")
	if err != nil {
		return err
	}
	port, err := getEnv("MAIL_PORT")
	if err != nil {
		return err
	}
	mail, err := getEnv("MAIL_NAME")
	if err != nil {
		return err
	}
	host, err := getEnv("MAIL_HOST")
	if err != nil {
		return err
	}
	app, err = getEnv("APP_NAME")
	if err != nil {
		return err
	}
	je := jsonEmail{user, password, mail, host, port}

	config = &je
	auth = LoginAuth(config.User, config.Password)
	return nil
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}

func (m Mail) sendMessage(message string, to, from mail.Address) error {
	if config.Host == "smtp.gmail.com" {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "smtp.gmail.com",
		}

		conn, err := tls.Dial("tcp", "smtp.gmail.com:465", tlsconfig)
		if err != nil {
			return err
		}

		c, err := smtp.NewClient(conn, "smtp.gmail.com")
		if err != nil {
			return err
		}

		// Auth
		if err = c.Auth(auth); err != nil {
			return err
		}

		// To && From
		if err = c.Mail(from.Address); err != nil {
			return err
		}

		if err = c.Rcpt(to.Address); err != nil {
			return err
		}

		// Data
		w, err := c.Data()
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		c.Quit()

		return nil
	} else {
		return smtp.SendMail(fmt.Sprintf("%s:%s", config.Host, config.Port), auth, config.Email, []string{to.Address}, []byte(message))
	}
}

func (m Mail) Confirm(dest, code string) error {
	if code == "" {
		return pkg.NewAppError("empty code", nil, http.StatusBadRequest)
	}
	if dest == "" {
		return pkg.NewAppError("empty address", nil, http.StatusBadRequest)
	}
	from := mail.Address{
		Name:    app,
		Address: config.Email,
	}
	to := mail.Address{Address: dest}
	subject := fmt.Sprintf("Validacion de correo electronico en %s", app)

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n"
	t, err := template.ParseFiles("web/confirm.html")
	if err != nil {
		return err
	}
	s := send{
		App:  app,
		Code: code,
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, s)
	if err != nil {
		return err
	}
	message += buf.String() + "\r\n"

	return m.sendMessage(message, to, from)
}
