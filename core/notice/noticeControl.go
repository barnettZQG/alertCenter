package notice

import ()
import (
	//"fmt"
	"time"
	"alertCenter/models"
	"github.com/astaxie/beego"
)

var DefaultTimeout = 3 * time.Second
var DefaultSendInterval = 5 * time.Second

func NoticControl(alert *models.Alert) {
	beego.Info("start alert:", alert.Mark)
	defer beego.Info("end alert:", alert.Mark)
	ch,ok := GetChannelByMark(alert.Fingerprint().String())
	if !ok{
		beego.Error("Can not find the channel.")
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

	var timer = time.NewTimer(DefaultSendInterval)

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
			timer.Reset(DefaultSendInterval)
		//fmt.Println("send email")
		}
	}

}