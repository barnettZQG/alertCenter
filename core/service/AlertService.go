package service

import (
	"alertCenter/core/db"
	"alertCenter/models"

	"github.com/astaxie/beego"

	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AlertService struct {
	Session *db.MongoSession
}

//GetAlertService  获取servcie
func GetAlertService(session *db.MongoSession) *AlertService {
	return &AlertService{
		Session: session,
	}
}

//GetAlertByLabels 获取报警根据labels
func (e *AlertService) GetAlertByLabels(alert *models.Alert) (result *models.Alert, err error) {
	//start := time.Now()
	//defer fmt.Println("cost:",time.Now().Sub(start))
	mark := alert.Fingerprint().String()
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil, fmt.Errorf("get alert collection error")
	}
	err = coll.Find(bson.M{"mark": mark}).Select(nil).One(&result)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		beego.Error("Get alert by Mark " + mark + " error." + err.Error())
		return nil, err
	}
	return
}

//Update 更新可变信息
func (e *AlertService) Update(alert *models.Alert) bool {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return false
	}
	err := coll.Update(bson.M{"mark": alert.Mark}, bson.M{
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
		beego.Error("Update the alert Error By Mark " + alert.Mark + "," + err.Error())
		return false
	}
	return true
}

// Save 存储报警
func (e *AlertService) Save(alert *models.Alert) bool {
	//return e.Session.Insert("Alert", alert)
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return false
	}
	_, err := coll.Upsert(bson.M{"mark": alert.Mark}, alert)
	if err == nil {
		return true
	}
	return false
}

//FindByUser 根据receiver的name或者id获取报警信息
func (e *AlertService) FindByUser(user string, pageSize int, page int) (alerts []*models.Alert, err error) {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil, fmt.Errorf("get alert collection error")
	}
	err = coll.Find(bson.M{"receiver.usernames": user, "ishandle": bson.M{"$ne": 2}}).Skip(pageSize * (page - 1)).Limit(pageSize).Select(nil).All(&alerts)
	return
}

//FindByID 根据mark获取报警
func (e *AlertService) FindByID(ID string) (alert *models.Alert, err error) {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil, nil
	}
	err = coll.Find(bson.M{"mark": ID}).One(&alert)
	if err != nil && err.Error() != mgo.ErrNotFound.Error() {
		beego.Debug("find alert by id faild." + err.Error())
	}
	return
}

//FindAll 获取全部报警
func (e *AlertService) FindAll(pageSize int, page int) (alerts []*models.Alert, err error) {
	coll := e.Session.GetCollection("Alert")
	if coll == nil {
		return nil, fmt.Errorf("get alert collection error")
	}
	err = coll.Find(bson.M{"ishandle": bson.M{"$ne": 2}}).Skip(pageSize * (page - 1)).Limit(pageSize).Select(nil).All(&alerts)
	return
}

//FindHistory 获取history通过alert
func (e *AlertService) FindHistory(alert *models.Alert) (history *models.AlertHistory, err error) {
	coll := e.Session.GetCollection("AlertHistory")
	if coll == nil {
		return nil, fmt.Errorf("get AlertHistory collection error")
	}
	err = coll.Find(bson.M{"mark": alert.Fingerprint().String(), "startsat": alert.StartsAt}).One(&history)
	if err != nil {
		beego.Error("find alerthistory by mark and startsAt faild." + err.Error())
	}
	return
}

//UpdateHistory 更新history时间信息
func (e *AlertService) UpdateHistory(history *models.AlertHistory) {
	coll := e.Session.GetCollection("AlertHistory")
	if coll == nil {
		return
	}
	err := coll.Update(bson.M{"mark": history.Mark, "startsat": history.StartsAt}, bson.M{
		"$set": bson.M{
			"endsat":   history.EndsAt,
			"startsat": history.StartsAt,
		},
	})
	if err != nil {
		beego.Error("update alerthistory by mark and startsAt faild." + err.Error())
	}
}
