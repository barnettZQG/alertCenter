package core

import (
	"crypto/tls"
	"time"

	"alertCenter/models"
	"alertCenter/util"

	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
)

type AlertSend interface {
	StartWork()
	StopWork()
	SendAlert(alert *models.Alert)
}

type MailAlertSend struct {
	mailChan chan *MailMessage
	stopChan chan bool
}

func (e *MailAlertSend) GetMailDialer() *gomail.Dialer {
	mailServer := beego.AppConfig.String("mailServer")
	mailPort, _ := beego.AppConfig.Int("mailPort")
	mailUser := beego.AppConfig.String("mailUser")
	mailPassword := beego.AppConfig.String("mailPassword")
	d := gomail.NewDialer(mailServer, mailPort, mailUser, mailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: mailPort == 465}
	return d
}

func (e *MailAlertSend) SendMail(message ...*gomail.Message) {
	d := e.GetMailDialer()
	d.DialAndSend(message...)
}

func (e *MailAlertSend) StartWork() error {
	beego.Info("mail send init start")
	defer beego.Info("mail send init over")
	mailCount, err := beego.AppConfig.Int("mailCount")
	if err != nil {
		beego.Error("mailCount's type is not int ." + err.Error())
		return err
	}
	mailReCount, err := beego.AppConfig.Int("mailReCount")
	if err != nil {
		beego.Error("mailReCount's type is not int ." + err.Error())
		return err
	}
	if e.mailChan == nil {
		e.mailChan = make(chan *MailMessage, mailCount)
	}
	if e.stopChan == nil {
		e.stopChan = make(chan bool)
	}
	go func() {
		d := e.GetMailDialer()
		var s gomail.SendCloser
		var err error
		open := false
		util.Info("mail work start success")
		for {
			select {
			case m, ok := <-e.mailChan:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						beego.Error("Get mail dial error." + err.Error())
					}
					open = true
				}
				if err := gomail.Send(s, m.message); err != nil {
					beego.Error("send mail message error." + err.Error())
					m.errCount++
					if m.errCount < mailReCount {
						//5秒后重试
						go func(m *MailMessage) {
							time.Sleep(time.Second * 5)
							e.mailChan <- m
						}(m)
					}
				}
			case stop := <-e.stopChan:
				if stop {
					goto exit
				}
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						panic(err)
					}
					open = false
				}
			}
		}
	exit:
		util.Info("mail work stop success")
	}()
	return nil
}
func (e *MailAlertSend) StopWork() {
	if e.stopChan != nil {
		e.stopChan <- true
		close(e.stopChan)
	}
	if e.mailChan != nil {
		close(e.mailChan)
	}
}

type MailMessage struct {
	message  *gomail.Message
	errCount int
	alert    *models.Alert
}

func (e *MailAlertSend) GetMessage(body string, subject string, receiver ...string) *MailMessage {
	m := gomail.NewMessage()
	m.SetHeader("From", beego.AppConfig.String("mailFrom"))
	m.SetHeader("To", receiver...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	return &MailMessage{
		message:  m,
		errCount: 0,
	}
}

func (e *MailAlertSend) GetMessageByAlert(alert *models.Alert) *MailMessage {

	m := e.GetMessage("", "", "")
	m.alert = alert
	return m
}

func (e *MailAlertSend) SendAlert(alert *models.Alert) {
	m := e.GetMessageByAlert(alert)
	e.mailChan <- m
}

// func TestSendMail() {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", beego.AppConfig.String("mailFrom"))
// 	m.SetHeader("To", "zengqingguo@goyoo.com")
// 	m.SetHeader("Subject", "Hello!")
// 	m.SetBody("text/plain", "Hello!")
// 	mailChan <- m
// 	stopChan <- true
// 	close(stopChan)
// 	close(mailChan)
// }
