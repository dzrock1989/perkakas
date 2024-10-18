package mail

import (
	"log"
	"strings"

	"github.com/dzrock1989/perkakas/configs"
	"gopkg.in/gomail.v2"
)

type mail struct {
	Subject    string
	BodyType   string
	Body       string
	Recipients []string
	CCs        []string
	BCCs       []string
}

func NewMail() mail {
	return mail{
		BodyType: "text/html",
	}
}

func (m *mail) PlainTextType() {
	m.BodyType = "text/plain"
}

func (m *mail) HTMLType() {
	m.BodyType = "text/html"
}

func (m *mail) SetSubject(subject string) {
	m.Subject = subject
}

func (m *mail) SetRecipient(rec string) {
	m.Recipients = append(m.Recipients, rec)
}

func (m *mail) SetCC(cc string) {
	m.CCs = append(m.CCs, cc)
}

func (m *mail) SetBCC(bcc string) {
	m.BCCs = append(m.BCCs, bcc)
}

func (m *mail) SetBodyType(bodyType string) {
	m.BodyType = bodyType
}

func (m *mail) SetBody(body string) {
	m.Body = body
}

func (m *mail) Print() {
	log.Println("========= KANG POS =========")
	log.Println("** Subject    \t:", m.Subject)
	log.Println("** Recipients \t:", strings.Join(m.Recipients, ", "))
	log.Println("** CCs        \t:", strings.Join(m.CCs, ", "))
	log.Println("** BCCs       \t:", strings.Join(m.BCCs, ", "))
	log.Println("** Type       \t:", m.BodyType)
	log.Println("** Body       \t:", m.Body)
	log.Println("============================")
}

func (m *mail) Send() (err error) {
	mm := gomail.NewMessage()
	mm.SetHeader("From", configs.Config.SMTP.Sender)
	mm.SetHeader("To", m.Recipients...)
	mm.SetHeader("Cc", m.CCs...)
	mm.SetHeader("Bcc", m.BCCs...)
	mm.SetHeader("Subject", m.Subject)
	mm.SetBody(m.BodyType, m.Body)

	d := gomail.NewDialer(
		configs.Config.SMTP.Host,
		configs.Config.SMTP.Port,
		configs.Config.SMTP.User,
		configs.Config.SMTP.Pass,
	)

	if err = d.DialAndSend(mm); err == nil {
		log.Println(
			"** Mail is sent for:",
			strings.Join(m.Recipients, ", "),
			strings.Join(m.CCs, ", "),
			strings.Join(m.BCCs, ", "),
		)
	}
	return
}
