package notice

import ()
import (
	//"fmt"
	"time"
	"alertCenter/models"
	"github.com/astaxie/beego"
)

var DefaultTimeout = 10 * time.Second

var SendMsgInterval_0 = time.Hour * 1
var SendMsgInterval_1 = time.Minute * 30
var SendMsgInterval_2 = time.Minute * 15
var SendMsgInterval_3 = time.Minute * 5

func init() {
	Lv0, err := time.ParseDuration(beego.AppConfig.String("sendMsgInterval_0"))
	if err == nil {
		SendMsgInterval_0 = Lv0
	} else {
		beego.Error(err)
	}
	Lv1, err := time.ParseDuration(beego.AppConfig.String("sendMsgInterval_1"))
	if err == nil {
		SendMsgInterval_1 = Lv1
	} else {
		beego.Error(err)
	}
	Lv2, err := time.ParseDuration(beego.AppConfig.String("sendMsgInterval_2"))
	if err == nil {
		SendMsgInterval_2 = Lv2
	} else {
		beego.Error(err)
	}
	Lv3, err := time.ParseDuration(beego.AppConfig.String("sendMsgInterval_3"))
	if err == nil {
		SendMsgInterval_3 = Lv3
	} else {
		beego.Error(err)
	}

}

func NoticControl(alert *models.Alert) {
	beego.Info("start alert:", alert.Mark)
	defer beego.Info("end alert:", alert.Mark)
	ch, err := GetChannelByMark(alert.Fingerprint().String())
	if err != nil {
		beego.Error(err)
		return
	}
	defer DeleteChanByMark(alert.Fingerprint().String())

	var timeout = time.NewTimer(DefaultTimeout)
	noNeedFlag := false

	NoNeedSend:
	for {
		select {
		case tmp := <-ch:
			if tmp.EndsAt.After(alert.StartsAt) {
				beego.Info("No need send this alert.")
			}
			noNeedFlag = true
		case <-timeout.C:
			break NoNeedSend
		}
	}

	if noNeedFlag {
		return
	}

	var timer = time.NewTimer(0 * time.Second)

	for {

		select {
		case al := <-ch:
			_ = alert
			if al.EndsAt.After(alert.StartsAt) {
				beego.Info("Alert has been fix")
			}
			return
		case <-timer.C:
			beego.Info("Sending email")
			//cacheServer["mail"].SendAlert(alert)
			timer.Reset(GetSendMsgInterval(alert.Level))
		//fmt.Println("send email")
		}
	}

}

func GetSendMsgInterval(level int) time.Duration {
	switch level{
	case 0:
		return SendMsgInterval_0
	case 1:
		return SendMsgInterval_1
	case 2:
		return SendMsgInterval_2
	case 3:
		return SendMsgInterval_3
	default:
		return SendMsgInterval_0
	}
}