package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/resend/resend-go/v3"
)

// Below we declare a new variable with the type embed.FS (embedded file system) to hold
// our email templates. This has a comment directive in the format `//go:embed <path>`
// IMMEDIATELY ABOVE it, which indicates to Go that we want to store the contents of the
// ./templates directory in the templateFS embedded file system variable.
// ↓↓↓

//go:embed "templates"
var templateFS embed.FS

// Define a Mailer struct which contains a mail.Dialer instance (used to connect to a
// SMTP server) and the sender information for your emails (the name and address you
// want the email to be from, such as "Alice Smith <alice@example.com>").
type Mailer struct {
	client *resend.Client
	sender string
}

func New(apiKey, sender string) Mailer {
	// Initialize a new resend mail client instance with the given apiKey server settings.
	client := resend.NewClient(apiKey)

	// Return a Mailer instance containing the dialer and sender information.
	return Mailer{
		client: client,
		sender: sender,
	}
}

// Define a Send() method on the Mailer type. This takes the receipient email address
// as the first param, the name of the file containing the templates, and any dynamic data for the
// templates as an interface{} param.
func (m Mailer) Send(recipient, templateFile string, data any) error {
	// Use the ParseFS() method to parse the required template file from the
	// embedded file system.
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

	// Build the Resend request payload
	params := &resend.SendEmailRequest{
		From:    m.sender,
		To:      []string{recipient},
		Subject: subject.String(),
		Text:    plainBody.String(),
		Html:    htmlBody.String(),
	}

	// Call the Send() method on the client.Emails, passing in the payload to send.

	// Try sending the email up to three times before aborting and returning the final
	// error. We sleep for 500 milliseconds between each attempt.
	for i := 1; i <= 3; i++ {
		_, err = m.client.Emails.Send(params)
		// If everything worked, return nil.
		if nil == err {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
