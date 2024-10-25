package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
)

var auth smtp.Auth

func main() {
	type emailaddress struct {
		name  string
		email string
	}
	// optional variables
	var msgfrom emailaddress
	var msgto []emailaddress
	// msgto := make([]emailaddress,1)
	// var cc emailaddress

	//
	// Required
	//
	to := []string{"craig.blueskyair@gmail.com", "craig@blueskyflying.com.au"}
	envelopeTo := "" // set later to create envelope data
	envelopeFrom := ""
	from := "dick@hammond.zone"
	subject := "gagf you cuntos" //technically email will work without a subject, but it's stupid

	// The from email address will be shown on the users email as who the email came from
	// but this can be overriden the the msgfrom: envelope variable so it includes a name
	// that shows up along with their email address.
	// the email address in from: and msgfrom: don't even have to match.

	//
	// Recommended
	//
	msgto = []emailaddress{
		{"Craigus", "craig@greenskyflying.com.au"},
		{"SITREP", "sitrep@sitrep"},
	}

	// Show to the user who the email was sent to. While it can match who it was actually sent to, there is
	// no requirement for it to match in any way.
	// To BCC, just

	//
	// Optional
	//
	// message envelope.
	// What is here doesn't affect who gets sent the email, but what shows up in the users email display.
	// It doesn't need to match either the required to or from variables above.
	// You can include name and email, or just email.

	//You can set the email address by itself, or name and email, but not just the name.
	//if both are set, it will get put in the header as  "craigus hammondus <craigus@hammondus.com>"

	// If message from isn't set here, it will default to the "from" variable in the Required section.
	// Setting it here allows you to put a name with the email address
	msgfrom.name = "Dick Breath"
	msgfrom.email = "dickbreath@hammond.zone"

	// cc = {"c.hammond@southernairlines.com.au>"}
	// bcc = "c.hammond@southernairlines.com.au"

	// Do some error checking before attempting to send the email
	// If msgfrom or msgto is used, check it's name & email, or just email, but not just name

	// if msgfrom.name != "" && msgfrom.email == "" {
	// 	msgfrom.name = ""
	// }
	// if msgfrom.name == "" && msgfrom.email == "" {
	// 	envelopeFrom = "From: " + from + "\r\n"
	// }

	if msgfrom.email == "" {
		envelopeFrom = "From: " + from + "\r\n"
	} else {
		envelopeFrom = "From: " + msgfrom.name + " <" + msgfrom.email + ">"
	}

	// if msgto.email == "" {
	// 	log.Println("setting To: to default")
	// 	envelopeTo = "To: " + strings.Join(to, ",") + "\r\n"
	// }
	// if msgto.name != "" && msgto.email != "" {
	// 	envelopeTo = "To: " + msgto.name + " <" + msgto.email + ">\r\n"
	// }

	envelopeTo = "To: "
	if len(msgto) != 0 {
		for _, v := range msgto {
			envelopeTo += fmt.Sprintf("%s <%s>,", v.name, v.email)
		}
		envelopeTo = strings.TrimSuffix(envelopeTo, ",")
	} else {
		envelopeTo += strings.Join(to, ",")
	}

	//create the message text
	msg := envelopeTo + "\r\n" // msg := "To: Cocko Breath <cockbreath@wankerfucker.com>\r\n"
	msg += envelopeFrom + "\r\n"
	msg += "Subject: " + subject + "\r\n"
	msg += "\r\n" +
		"This is the email bodyyyyyyy.\r\n"

	fmt.Println(msg)
	mymsg := []byte(msg)

	// msg := []byte("To: bdstudy <bdstud@hammond.zone>\r\n" +
	// 	// "Subject: gagffff\r\n" +
	// 	"\r\n" +
	// 	"This is the email body.\r\n")

	testes := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, testes, from, to, mymsg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

	emailTo := []string{"c.hammond@southernairlines.com.au", "craig.blueskyair@gmail.com"}
	// emailFrom := "craig@blueskyflying.com.au"

	auth = smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)
	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Craigus Hammondus",
		URL:  "http://app.southernairlines.com.au",
	}
	r := NewRequest(emailTo, "New Southern App", "Hello, World!")
	// err := r.ParseTemplate("template.html", templateData)
	if err := r.ParseTemplate("template.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	}

}

// Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	// addr := "email-smtp.ap-southeast-2.amazonaws.com:587"
	addr := smtpServer + ":" + smtpPort

	if err := smtp.SendMail(addr, auth, "craig@blueskyflying.com.au", r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
