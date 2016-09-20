package notice

import (
	//"fmt"
	"alertCenter/models"
	"time"

	"github.com/astaxie/beego"
)

var IgnoreSend = 10 * time.Second
var StopSend = 30 * time.Minute

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

	ignoreSend, err := time.ParseDuration(beego.AppConfig.String("ignoreSend"))
	if err == nil {
		IgnoreSend = ignoreSend
	} else {
		beego.Error(err)
	}

	stopSend, err := time.ParseDuration(beego.AppConfig.String("stopSend"))
	if err == nil {
		StopSend = stopSend
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

	var timeout = time.NewTimer(IgnoreSend)
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
	var stopSend = time.NewTimer(StopSend)

	for {

		select {
		case al := <-ch:
			_ = alert
			if al.EndsAt.After(alert.StartsAt) {
				beego.Info("Alert has been fix")
				return
			}
			stopSend.Reset(StopSend)
		case <-stopSend:
			beego.Info("Have not get this alert for long time. Stop sending email.")
			return
		case <-timer.C:
		//if beego.AppConfig.String("runmode") != "dev" {
			for _, server := range cacheServer {
				if server != nil {
					server.SendAlert(alert)
				}
			}
		//}

			timer.Reset(GetSendMsgInterval(alert.Level))
		}
	}

}

func GetSendMsgInterval(level int) time.Duration {
	switch level {
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
