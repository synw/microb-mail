package mail

import (
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/synw/microb-http/csrf"
	"github.com/synw/microb-http/types"
	config "github.com/synw/microb-mail/conf"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/http"
	"os"
)

var Conf *config.Conf
var Dev bool
var mailTemplate *template.Template

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

func Init(dev bool) *terr.Trace {
	Dev = dev
	conf, tr := config.GetConf()
	if tr != nil {
		events.Error("mail", "Can not get mail service config", tr, "fatal")
		return tr
	}
	Conf = conf
	return nil
}

func ParseTemplate() {
	path, _ := os.Getwd()
	path = path + "/templates/*"
	mailTemplate = template.Must(template.ParseGlob(path))
}

func ServeMailForm(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w = httpResponseWriter{w, &status}
	token, tr := csrf.GetToken()
	if tr != nil {
		tr := terr.Add("mail.ServeMailForm", errors.New("Error serving mail form"), tr)
		events.Error("mail", "Can not serve mail form", tr)
	}
	content := ""
	page := &types.Page{"", "/mail", "Email", template.HTML(content), &types.Conn{}, "", token}
	err := mailTemplate.ExecuteTemplate(w, "mail_form.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering template: " + err.Error()
		err := errors.New("Can not render template")
		tr := terr.New("mail.ServeMailForm", err)
		events.Error("mail", msg, tr)
	}
}

func ProcessMailForm(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	email := r.FormValue("email")
	subject := r.FormValue("subject")
	msg := r.FormValue("content")
	msg = sanitizeInput(msg)
	email = sanitizeInput(email)
	subject = sanitizeInput(subject)
	tr := csrf.VerifyToken(token)
	if tr != nil {
		events.Error("mail", "Can not process mail form", tr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	// send mail
	if Dev == false {
		sendMail(email, subject, msg)
	}
	// respond
	//status := http.StatusOK
	//w = httpResponseWriter{w, &status}
	http.Redirect(w, r, "/mail/ok", http.StatusMovedPermanently)

}

func sendMail(from string, subject string, msg string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", Conf.To)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", msg)
	d := gomail.NewDialer(Conf.Host, Conf.Port, Conf.User, Conf.Password)
	//d := gomail.NewDialer("localhost", 25, "", "")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	events.Info("mail", "Sending mail from "+from)
}

func sanitizeInput(input string) string {
	p := bluemonday.NewPolicy()
	output := p.Sanitize(input)
	return output
}
