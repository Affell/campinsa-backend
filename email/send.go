package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/kataras/golog"
)

func (email Structure) Send(config Config) (success bool) {

	from := mail.Address{
		Name:    config.From,
		Address: config.User,
	}

	to := mail.Address{
		Name:    "",
		Address: email.To,
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = html.EscapeString(email.Subject)
	headers["MIME-version"] = "1.0;"
	headers["Content-Type"] = `text/html; charset="utf-8"`

	var message string

	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	buff := new(bytes.Buffer)
	err := config.Template.ExecuteTemplate(buff, "email", email.Vars)
	if err != nil {
		fmt.Println(buff.String())
		golog.Error(err.Error())
		return
	}

	message += "\r\n" + buff.String()

	auth := login(config.User, config.Password)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.Host,
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(config.Host, config.Port))
	if err != nil {
		golog.Error("net.Dial('tcp', net.JoinHostPort(config.Host, config.Port)) when sending mail")
		return
	}

	defer conn.Close()

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		golog.Error("smtp.NewClient(conn, config.Host) when sending mail")
		return
	}

	defer client.Quit()

	if err = client.StartTLS(tlsconfig); err != nil {
		golog.Error("client.StartTLS(tlsconfig) when sending mail")
		return
	}

	if err = client.Auth(auth); err != nil {
		golog.Error("client.Auth(auth) when sending mail")
		return
	}

	if err = client.Mail(from.Address); err != nil {
		golog.Error("client.Mail(from.Address) when sending mail")
		return
	}

	if err = client.Rcpt(to.Address); err != nil {
		golog.Error("client.Rcpt when sending mail")
		return
	}

	w, err := client.Data()
	if err != nil {
		golog.Error("Client.Data() when sending mail")
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		golog.Error("writing message in reader")
		return
	}

	err = w.Close()
	if err != nil {
		golog.Error("during closing connection to the email server")
		return
	}

	return true
}
