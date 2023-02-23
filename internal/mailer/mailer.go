package mailer

import (
	"bytes"
	"embed"
	"github.com/go-mail/mail/v2"
	"html/template"
	"time"
)

// Below we declare a new variable with the type embed.FS (embedded file system) to hold
// our email templates. This has a comment directive in the format `//go:embed <path>`
// IMMEDIATELY ABOVE it, which indicates to Go that we want to store the contents of the
// ./templates directory in the templateFS embedded file system variable.
// ↓↓↓
//
//go:embed "templates"
var templateFS embed.FS

// define a Mailer instance which contains a mail.Dialer instance (used to connect to the SMTP server)
// and instance of sender information for emails (the name and address which the email will be from)
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// constructor for our Mailer
func New(host string, port int, username string, password string, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// define Send() method on the Mailer type. it takes recipient email, name of the file containing our templates
// and any dynamic data for templates as any
func (m Mailer) Send(recipient string, templateFile string, data any) error {
	// using ParseFS() method in order to parse required file template from the EFS(embedded file system)
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	// And likewise with the "htmlBody" template.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// create a mail.NewMessage() instance, in order to create new messages
	msg := mail.NewMessage()
	// Then we use the SetHeader() method to set the email recipient, sender and subject
	// headers, the SetBody() method to set the plain-text body, and the AddAlternative()
	// method to set the HTML body. It's important to note that AddAlternative() should
	// always be called *after* SetBody().
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())
	msg.Attach("internal / mailer / img / welcome.png")

	// then call the DialAndSend() method, pass in message. it opens a connection to SMTP server
	// sends the message and closes connection. if there is a timeout, it will return it
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
