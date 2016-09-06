package core

import (
	"time"

	"github.com/barnettZQG/alertCenter/models"
	"github.com/barnettZQG/alertCenter/util"
)

//HandleMessage 处理alertmanager回来的数据
func HandleMessage(message *models.AlertReceive) {
	session := GetMongoSession()
	defer session.Close()
	alertService := &AlertService{
		Session: session,
	}
	ok := SaveMessage(message, session)
	if !ok {
		util.Error("save a message fail,message receiver:" + message.Receiver)
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
			alert.Receiver = GetReceiverByTeam(message.Receiver)
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
	session := GetMongoSession()
	defer session.Close()
	alertService := &AlertService{
		Session: session,
	}
	for _, alert := range alerts {
		old := alertService.GetAlertByLabels(alert)
		if old != nil {
			old.AlertCount = old.AlertCount + 1
			alert.UpdatedAt = time.Now()
			old = old.Merge(alert)
			if !old.EndsAt.IsZero() {
				old.IsHandle = 2
				old.HandleDate = time.Now()
				old.HandleMessage = "报警已自动恢复"
			}
			old.UpdatedAt = time.Now()
			alertService.Update(old)
		} else {
			alert.AlertCount = 1
			alert.IsHandle = 0
			alert.Mark = alert.Fingerprint().String()
			alert.Receiver = GetReceiver(alert.Labels)
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
			alert.UpdatedAt = now
			alertService.Save(alert)
		}
	}
}

func SaveMessage(message *models.AlertReceive, session *MongoSession) bool {
	ok := session.Insert("AlertReceive", message)
	return ok
}
