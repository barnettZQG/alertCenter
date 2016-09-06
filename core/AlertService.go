package core

import (
	"github.com/barnettZQG/alertCenter/models"
	"github.com/barnettZQG/alertCenter/util"
	"gopkg.in/mgo.v2/bson"
)

type AlertService struct {
	Session *MongoSession
}

//GetAlertService  获取servcie
func GetAlertService(session *MongoSession) *AlertService {
	return &AlertService{
		Session: session,
	}
}

//GetAlertByLabels 获取报警根据labels
func (e *AlertService) GetAlertByLabels(alert *models.Alert) (result *models.Alert) {
	mark := alert.Fingerprint().String()
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil
	}
	err := coll.Find(bson.M{"mark": mark}).Select(nil).One(&result)
	if err != nil {
		util.Debug("Get alert by Mark " + mark + " error." + err.Error())
		return nil
	}
	return
}

//Update 更新可变信息
func (e *AlertService) Update(alert *models.Alert) bool {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return false
	}
	err := coll.Update(bson.M{"mark": alert.Mark, "ishandle": 0}, bson.M{
		"$set": bson.M{
			"alertcount":    alert.AlertCount,
			"ishandle":      alert.IsHandle,
			"handledate":    alert.HandleDate,
			"handlemessage": alert.HandleMessage,
			"endsat":        alert.EndsAt,
			"startsat":      alert.StartsAt,
			"updatedat":     alert.UpdatedAt,
		},
	})
	if err != nil {
		util.Debug("Update the alert Error By Mark " + alert.Mark + "," + err.Error())
		return false
	}
	return true
}

// Save 存储报警
func (e *AlertService) Save(alert *models.Alert) bool {
	return e.Session.Insert("Alert", alert)
}

//FindByUser 根据receiver的name或者id获取报警信息
func (e *AlertService) FindByUser(user string) (alerts []*models.Alert) {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil
	}
	coll.Find(bson.M{"receiver.name": user, "ishandle": 0}).Select(nil).All(&alerts)
	if alerts == nil || len(alerts) == 0 {
		coll.Find(bson.M{"receiver.id": user, "ishandle": 0}).Select(nil).All(&alerts)
	}
	return
}

//FindByID 根据ID获取报警
func (e *AlertService) FindByID(ID string) (alert *models.Alert) {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil
	}
	err := coll.Find(bson.M{"id": ID}).One(&alert)
	if err != nil {
		util.Debug("find alert by id faild." + err.Error())
	}
	return
}

//FindAll 获取全部报警
func (e *AlertService) FindAll() (alerts []*models.Alert) {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil
	}
	coll.Find(bson.M{"ishandle": bson.M{"$ne": 2}}).Select(nil).All(&alerts)
	return
}
