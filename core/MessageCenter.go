package core

import (
	"time"

	"github.com/astaxie/beego"

	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/core/user"
	"alertCenter/models"
)

//HandleMessage 处理alertmanager回来的数据
func HandleMessage(message *models.AlertReceive) {
	session := db.GetMongoSession()
	defer session.Close()
	alertService := &service.AlertService{
		Session: session,
	}
	ok := SaveMessage(message, session)
	if !ok {
		beego.Error("save a message fail,message receiver:" + message.Receiver)
	}
	for _, alert := range message.Alerts {
		old := alertService.GetAlertByLabels(&alert)
		if old != nil {
			old.AlertCount = old.AlertCount + 1
			old = old.Merge(&alert)
			if !old.EndsAt.IsZero() {
				old.IsHandle = 2
				old.HandleDate = time.Now()
				old.HandleMessage = "报警已自动恢复"
			}
			alertService.Update(old)
		} else {
			alert.AlertCount = 1
			alert.IsHandle = 0
			alert.Mark = alert.Fingerprint().String()
			alert.Receiver = user.GetReceiverByTeam(message.Receiver)
			now := time.Now()
			// Ensure StartsAt is set.
			if alert.StartsAt.IsZero() {
				alert.StartsAt = now
			}
			if !alert.EndsAt.IsZero() {
				alert.IsHandle = 2
				alert.HandleDate = time.Now()
				alert.HandleMessage = "报警已自动恢复"
			}
			alertService.Save(&alert)
		}
	}

}

//HandleAlerts 处理prometheus回来的数据
func HandleAlerts(alerts []*models.Alert) {
	session := db.GetMongoSession()
	defer session.Close()
	alertService := &service.AlertService{
		Session: session,
	}
	for _, alert := range alerts {
		old := alertService.GetAlertByLabels(alert)
		if old != nil && old.EndsAt.IsZero() {
			old.AlertCount = old.AlertCount + 1
			alert.UpdatedAt = time.Now()
			old = old.Merge(alert)
			//old已更新时间信息
			if !old.EndsAt.IsZero() {
				old.IsHandle = 2
				old.HandleDate = time.Now()
				old.HandleMessage = "报警已自动恢复"
				SaveHistory(alertService, old)
			}
			old.UpdatedAt = time.Now()
			Notice(old)
			alertService.Update(old)
		} else if old != nil && !old.EndsAt.IsZero() {
			//此报警曾出现过并已结束
			if alert.StartsAt.After(old.EndsAt) {
				//报警开始时间在原报警之后，我们认为这是新报警
				//old更新状态信息
				old = old.Reset(alert)
				if old.IsHandle == 2 {
					SaveHistory(alertService, old)
				}
				Notice(old)
				alertService.Update(old)
			} else if alert.StartsAt.Before(old.EndsAt) && alert.EndsAt.After(old.EndsAt) {
				// 新的结束时间
				history := alertService.FindHistory(old)
				old.EndsAt = alert.EndsAt
				history.EndsAt = alert.EndsAt
				alertService.Update(old)
				alertService.UpdateHistory(history)
			}
		} else {
			//曾经没出现过的报警
			SaveAlert(alertService, alert)
		}
	}
}

//Notice 发送报警通知信息
func Notice(alert *models.Alert) {
	// if users, ok := CheckRules(alert); ok {
	// 	alert.Receiver.UserNames = users
	// 	mark := alert.Fingerprint().String()

	// 	ch, err := notice.GetChannelByMark(mark)
	// 	if err == nil && ch != nil {
	// 		ch <- alert
	// 	} else {
	// 		err := notice.CreateChanByMark(alert.Fingerprint().String())
	// 		if err != nil {
	// 			beego.Error(err)
	// 		}
	// 		go notice.NoticControl(alert)
	// 	}
	// }
}

//CheckRules 检验是否为用户忽略的报警
func CheckRules(alert *models.Alert) ([]string, bool) {
	if alert != nil && alert.Receiver != nil {
		userNames := alert.Receiver.UserNames
		//没有接收对象，不需要发送
		if userNames == nil || len(userNames) < 1 {
			return nil, false
		}
		relation := user.Relation{}
		var users []string
		session := db.GetMongoSession()
		defer session.Close()
		ruleService := &service.IgnoreRuleService{
			Session: session,
		}
		for _, userName := range userNames {
			user := relation.GetUserByName(userName)
			var ignore bool
			if user != nil {
				rules := ruleService.FindRuleByUser(user.ID)
				if rules != nil && len(rules) > 0 {
					for _, rule := range rules {
						if alert.Labels.Contains(rule.Labels) {
							ignore = true
						}
					}
				}
			} else {

			}
			if !ignore {
				users = append(users, user.Name)
			}
		}
		if len(users) == 0 {
			return nil, false
		}
		return users, true
	}
	return nil, false
}

//SaveHistory 存快照纪录
func SaveHistory(alertService *service.AlertService, alert *models.Alert) {
	history := &models.AlertHistory{
		Mark:     alert.Fingerprint().String(),
		StartsAt: alert.StartsAt,
		EndsAt:   alert.EndsAt,
		Message:  string(alert.Annotations.LabelSet["description"]),
	}
	alertService.Session.Insert("AlertHistory", history)
}

//SaveAlert 保存alert信息
func SaveAlert(alertService *service.AlertService, alert *models.Alert) {
	alert.AlertCount = 1
	alert.IsHandle = 0
	alert.Mark = alert.Fingerprint().String()
	alert.Receiver = user.GetReceiver(alert.Labels)
	now := time.Now()
	// Ensure StartsAt is set.
	if alert.StartsAt.IsZero() {
		alert.StartsAt = now
	}
	alert.UpdatedAt = now
	if !alert.EndsAt.IsZero() {
		alert.IsHandle = 2
		alert.HandleDate = time.Now()
		alert.HandleMessage = "报警已自动恢复"
		SaveHistory(alertService, alert)
	}
	Notice(alert)
	alertService.Save(alert)
}

//SaveMessage 储存alertmanager的消息
func SaveMessage(message *models.AlertReceive, session *db.MongoSession) bool {
	ok := session.Insert("AlertReceive", message)
	return ok
}
