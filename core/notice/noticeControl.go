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

	timeout := make(chan bool)
	go func() {
		time.Sleep(DefaultTimeout)
		timeout <- true
	}()
	noNeedFlag := false

	ch := GetChannelByMark(alert.Fingerprint().String())
	defer DeleteChanByMark(alert.Fingerprint().String())
	NoNeedSend:
	for {
		select {
		case tmp := <-ch:
			if tmp.EndsAt.After(alert.StartsAt) {
				beego.Info("No need send this alert.")
			}
			noNeedFlag = true
		case <-timeout:
		//fmt.Println("timeout 1")
			break NoNeedSend
		}
	}

	if noNeedFlag {
		return
	}

	//tout := make(chan bool)
	//go func() {
	//	time.Sleep(DefaultSendInterval)
	//	tout <- true
	//}()
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
			timer.Reset(DefaultSendInterval)
		//fmt.Println("send email")
		}
	}

}