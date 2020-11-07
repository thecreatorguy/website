package message

import (
	"crypto/tls"
	"net/http"
	"net/smtp"
	"os"
	"website/internal/app/response"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)


func sendMessageEndpoint(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.Write400(w)
		return
	}

	e := email.NewEmail()
	e.From = "Website Questions <questions.for.tim@itstimjohnson.com>"
	e.To = []string{"tim@itstimjohnson.com"}
	e.Cc = []string{r.Form.Get("email")}
	e.Subject = "A Message From " + r.Form.Get("name")
	e.Text = []byte(r.Form.Get("message"))
	err = e.SendWithTLS(
		os.Getenv("EMAIL_HOST") + ":" + os.Getenv("EMAIL_PORT"),
		smtp.PlainAuth("", "questions.for.tim@itstimjohnson.com", os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")),
		&tls.Config{ServerName: os.Getenv("EMAIL_HOST")},
	)
	if err != nil {
		logrus.WithError(err).Error("Failed sending email")
		response.Write500(w)
		return
	}

	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}