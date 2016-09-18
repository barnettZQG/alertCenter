package service

import (
	"alertCenter/core/db"
	"alertCenter/models"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type GlobalConfigService struct {
	Session *db.MongoSession
}

//GetConfig 获取全局配置
func (e *GlobalConfigService) GetConfig(name string) (config *models.GlobalConfig) {
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		err := coll.Find(bson.M{"name": name}).Select(nil).One(&config)
		if err != nil {
			beego.Error("get GlobalConfig with name " + name + " error " + err.Error())
		}
	}
	return
}

//Update 更新全局配置
func (e *GlobalConfigService) Update(config *models.GlobalConfig) bool {
	if config == nil {
		return false
	}
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		err := coll.Update(bson.M{"name": config.Name}, bson.M{"$set": bson.M{"value": config.Value}})
		if err != nil {
			beego.Error("update GlobalConfig with name " + config.Name + " error " + err.Error())
			return false
		}
		return true
	}
	return false
}
