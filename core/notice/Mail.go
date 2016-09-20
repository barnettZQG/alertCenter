package notice

import (
	"crypto/tls"
	"path/filepath"
	"time"

	"alertCenter/core/user"
	"alertCenter/models"

	"io/ioutil"

	"strings"

	"strconv"

	"github.com/astaxie/beego"
	"gopkg.in/gomail.v2"
)

type MailNoticeServer struct {
	mailChan chan *MailMessage
	stopChan chan bool
}

//GetMailDialer 获取邮箱服务器代理
func (e *MailNoticeServer) GetMailDialer() *gomail.Dialer {
	mailServer := beego.AppConfig.String("mailServer")
	mailPort, _ := beego.AppConfig.Int("mailPort")
	mailUser := beego.AppConfig.String("mailUser")
	mailPassword := beego.AppConfig.String("mailPassword")
	d := gomail.NewDialer(mailServer, mailPort, mailUser, mailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: mailPort == 465}
	return d
}

//SendMail 发送邮件
func (e *MailNoticeServer) SendMail(message ...*gomail.Message) {
	d := e.GetMailDialer()
	d.DialAndSend(message...)
}

//StartWork 开始工作
func (e *MailNoticeServer) StartWork() error {
	beego.Info("mail notice server init start")
	defer beego.Info("mail notice server init over")
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
						beego.Debug("mail errCount:", m.errCount)
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
		beego.Info("mail work stop success")
	}()
	beego.Info("mail notice server start success")
	return nil
}

//StopWork 结束工作
func (e *MailNoticeServer) StopWork() error {
	if e.stopChan != nil {
		e.stopChan <- true
		close(e.stopChan)
	}
	if e.mailChan != nil {
		close(e.mailChan)
	}
	return nil
}

type MailMessage struct {
	message  *gomail.Message
	errCount int
	alert    *models.Alert
}

//GetMessage 构建邮件消息
func (e *MailNoticeServer) GetMessage(body string, subject string, receiver ...string) *MailMessage {
	m := gomail.NewMessage()
	m.SetHeader("From", beego.AppConfig.String("mailFrom"))
	m.SetHeader("To", receiver...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return &MailMessage{
		message:  m,
		errCount: 0,
	}
}

//GetMessageByAlert 通过alert获取邮件消息
func (e *MailNoticeServer) GetMessageByAlert(alert *models.Alert) (messages []*MailMessage) {
	userNames := alert.Receiver.UserNames
	relation := user.Relation{}
	for _, userName := range userNames {
		user := relation.GetUserByName(userName)
		if user != nil && user.Mail != "" {
			m := e.GetMessage(e.GetBody(alert), string(alert.Labels.LabelSet["alertname"])+"("+strconv.Itoa(alert.AlertCount)+")", user.Mail)
			m.alert = alert
			messages = append(messages, m)
		} else {
			beego.Debug("send mail to " + userName + ",user is not exit")
		}
	}
	return
}

//GetBody 创建邮件内容
func (e *MailNoticeServer) GetBody(alert *models.Alert) string {
	path, _ := filepath.Abs("views/mail.html")
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		beego.Error("get mail template file error." + err.Error())
	}
	mail := string(buffer)
	mail = strings.Replace(mail, "[TITLE]", string(alert.Labels.LabelSet["alertname"]), -1)
	mail = strings.Replace(mail, "[URL]", beego.AppConfig.String("url")+"/alertsCurrent", -1)
	mail = strings.Replace(mail, "[DESCRIPTION]", string(alert.Annotations.LabelSet["description"]), -1)
	return mail
}

//SendAlert 发送报警邮件实现
func (e *MailNoticeServer) SendAlert(alert *models.Alert) error {
	//beego.Debug("start SendAlert ")
	messages := e.GetMessageByAlert(alert)
	//beego.Debug("mail count:" + strconv.Itoa(len(messages)))
	if len(messages) > 0 {
		for _, m := range messages {
			e.mailChan <- m
		}
	}
	return nil
}
