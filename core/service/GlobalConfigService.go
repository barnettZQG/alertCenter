package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

type GlobalConfigService struct {
	Session *db.MongoSession
}

func (e *GlobalConfigService) Init() error {
	if config := e.GetConfig("noticeOn"); config == nil {
		e.Session.Insert("GlobalConfig", &models.GlobalConfig{
			Name:    "noticeOn",
			Value:   true,
			AddTime: time.Now(),
		})
	}
	return nil
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

//GetConfigA 获取全局配置
func (e *GlobalConfigService) GetConfigA(name string, value interface{}) (config *models.GlobalConfig) {
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		err := coll.Find(bson.M{"name": name, "value": value}).Select(nil).One(&config)
		if err != nil {
			beego.Error("get GlobalConfig with name " + name + " error " + err.Error())
		}
	}
	return
}

//GetAllConfig 获取同名全部配置
func (e *GlobalConfigService) GetAllConfig(name string) (configs []*models.GlobalConfig) {
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		err := coll.Find(bson.M{"name": name}).Select(nil).All(&configs)
		if err != nil {
			beego.Error("get GlobalConfig with name " + name + " error " + err.Error())
		}
	}
	return
}

//CheckExist 判断指定变量是否存在
func (e *GlobalConfigService) CheckExist(name string, value interface{}) bool {
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		var config models.GlobalConfig
		err := coll.Find(bson.M{"name": name, "value": value}).Select(nil).One(&config)
		if err == nil && &config != nil {
			return true
		}
	}
	return false
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

//DeleteByID 删除配置通过id
func (e *GlobalConfigService) DeleteByID(id string) bool {
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		err := coll.RemoveId(bson.ObjectIdHex(id))
		if err != nil {
			beego.Error("remove GlobalConfig by ID(" + id + ") error ." + err.Error())
			return false
		}
		return true
	}
	return false
}
