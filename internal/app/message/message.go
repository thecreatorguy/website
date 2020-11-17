// Package message adds a route that handles sending me (Tim) a message from the "contact me" page
package message

import (
	"encoding/json"
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

	// Check recaptcha
	client := &http.Client{}
	res, err := client.PostForm("https://www.google.com/recaptcha/api/siteverify", map[string][]string{
		"secret": {os.Getenv("CAPTCHA_SECRET")},
		"response": {r.Form.Get("g-recaptcha-response")},
	})
	if err != nil {
		response.Write500(w)
		return
	}
	defer res.Body.Close()
	var decoded map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		response.Write500(w)
		return
	}
	if decoded["success"].(bool) != true {
		logrus.WithField("response", decoded).Warn("Captcha returned a failure")
		response.Write403(w)
		return
	}

	// Send email message to me
	e := email.NewEmail()
	e.From = "Website Questions <questions.for.tim@itstimjohnson.com>"
	e.To = []string{"tim@itstimjohnson.com"}
	e.Cc = []string{r.Form.Get("email")}
	e.Subject = "A Message From " + r.Form.Get("name")
	e.Text = []byte(r.Form.Get("message"))
	err = e.Send(
		os.Getenv("EMAIL_HOST") + ":" + os.Getenv("EMAIL_PORT"),
		smtp.PlainAuth("", os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")),
	)
	if err != nil {
		logrus.WithError(err).Error("Failed sending email")
		response.Write500(w)
		return
	}

	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}