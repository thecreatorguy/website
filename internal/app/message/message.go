// Package message adds a route that handles sending me (Tim) a message from the "contact me" page
package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"website/internal/app/response"
)


func sendMessageEndpoint(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.Write400(w)
		return
	}

	// Check recaptcha
	// var buf bytes.Buffer
	// json.NewEncoder(&buf).Encode(map[string]interface{}{
		
	// })
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
		fmt.Println(decoded)
		response.Write403(w)
	}

	// Send email message to me
	// e := email.NewEmail()
	// e.From = "Website Questions <questions.for.tim@itstimjohnson.com>"
	// e.To = []string{"tim@itstimjohnson.com"}
	// e.Cc = []string{r.Form.Get("email")}
	// e.Subject = "A Message From " + r.Form.Get("name")
	// e.Text = []byte(r.Form.Get("message"))
	// err = e.SendWithTLS(
	// 	os.Getenv("EMAIL_HOST") + ":" + os.Getenv("EMAIL_PORT"),
	// 	smtp.PlainAuth("", "questions.for.tim@itstimjohnson.com", os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST")),
	// 	&tls.Config{ServerName: os.Getenv("EMAIL_HOST")},
	// )
	// if err != nil {
	// 	logrus.WithError(err).Error("Failed sending email")
	// 	response.Write500(w)
	// 	return
	// }

	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}