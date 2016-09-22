package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GlobalConfigService struct {
	Session *db.MongoSession
}

var globalConfigMap map[string]*models.GlobalConfig
var globalConfigs []*models.GlobalConfig

//RefreshGlobalCnfig 更新全局配置缓存
func (e *GlobalConfigService) RefreshGlobalCnfig() {
	beego.Debug("Start refresh global configs.")
	if e.Session == nil {
		return
	}
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
	} else {
		var configs []*models.GlobalConfig
		err := coll.Find(nil).Select(nil).All(&configs)
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			beego.Error("get GlobalConfigs error " + err.Error())
		}
		if configs != nil {
			globalConfigs = configs
			globalConfigTmp := make(map[string]*models.GlobalConfig, 0)
			for _, config := range configs {
				globalConfigTmp[config.Name] = config
			}
			globalConfigMap = globalConfigTmp
			beego.Debug("refresh global configs success.size:", len(globalConfigMap))
		}
	}

}
func (e *GlobalConfigService) Init() error {
	globalConfigMap = make(map[string]*models.GlobalConfig, 0)
	globalConfigs = make([]*models.GlobalConfig, 0)
	if config, _ := e.GetConfig("noticeOn"); config == nil {
		e.Insert(&models.GlobalConfig{
			ID:      bson.NewObjectId(),
			Name:    "noticeOn",
			Value:   false,
			AddTime: time.Now(),
		})
	}
	e.RefreshGlobalCnfig()
	return nil
}

//GetConfig 获取全局配置,重复名称的配置不适合此方法
func (e *GlobalConfigService) GetConfig(name string) (config *models.GlobalConfig, err error) {
	if config := globalConfigMap[name]; config != nil {
		return config, nil
	}
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
		err = fmt.Errorf("get GlobalConfig collection error")
	} else {
		err := coll.Find(bson.M{"name": name}).Select(nil).One(&config)
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			beego.Error("get GlobalConfig with name " + name + " error " + err.Error())
		}
	}
	return
}

//GetConfigA 获取全局配置
func (e *GlobalConfigService) GetConfigA(name string, value interface{}) (config *models.GlobalConfig, err error) {
	if globalConfigs != nil {
		for _, config := range globalConfigs {
			if config.Name == name && config.Value == value {
				return config, nil
			}
		}
	}
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
		err = fmt.Errorf("get GlobalConfig collection error")
	} else {
		err = coll.Find(bson.M{"name": name, "value": value}).Select(nil).One(&config)
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			beego.Error("get GlobalConfig with name " + name + " error " + err.Error())
		}
	}
	return
}

//GetAllConfig 获取同名全部配置
func (e *GlobalConfigService) GetAllConfig(name string) (configs []*models.GlobalConfig, err error) {
	if globalConfigs != nil {
		for _, config := range globalConfigs {
			if config.Name == name {
				configs = append(configs, config)
			}
		}
		return
	}
	coll := e.Session.GetCollection("GlobalConfig")
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
		err = fmt.Errorf("get GlobalConfig collection error")
	} else {
		err = coll.Find(bson.M{"name": name}).Select(nil).All(&configs)
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			beego.Error("get GlobalConfig with name " + name + " error " + err.Error())
		}
	}
	return
}

//CheckExist 判断指定变量是否存在, 适用于value类型为string的
func (e *GlobalConfigService) CheckExist(name string, value interface{}) (bool, error) {
	if globalConfigs != nil {
		for _, config := range globalConfigs {
			//beego.Debug(config.Name, name, config.Value, value)
			if config.Name == name && config.Value.(string) == value.(string) {
				return true, nil
			}
		}
		return false, nil
	}
	coll := e.Session.GetCollection("GlobalConfig")
	var err error
	if coll == nil {
		beego.Error("get GlobalConfig collection error when init NoticeOn ")
		err = fmt.Errorf("get GlobalConfig collection error")
	} else {
		var config *models.GlobalConfig
		err = coll.Find(bson.M{"name": name, "value": value}).Select(nil).One(&config)
		if err != nil && err.Error() != mgo.ErrNotFound.Error() {
			beego.Error("get config error," + err.Error())
			return false, err
		}
		if err != nil && err.Error() == mgo.ErrNotFound.Error() {
			return false, nil
		}
		if &config != nil {
			return true, nil
		}
	}
	return false, err
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
		if globalConfigMap != nil {
			globalConfigMap[config.Name] = config
		}
		if globalConfigs != nil {
			for _, c := range globalConfigs {
				if c.Name == config.Name {
					c.Value = config.Value
				}
			}
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
		var co *models.GlobalConfig
		var cos []*models.GlobalConfig
		if globalConfigs != nil {
			for _, c := range globalConfigs {
				if c.ID.Hex() == id {
					co = c
				} else {
					cos = append(cos, c)
				}
			}
			globalConfigs = cos
			beego.Debug("delete config id:", id, "globalConfigs length :", len(globalConfigs))
		}
		if globalConfigMap != nil && co != nil {
			delete(globalConfigMap, co.Name)
		}
		return true
	}
	return false
}

//Insert 添加全局配置
func (e *GlobalConfigService) Insert(config *models.GlobalConfig) bool {
	if e.Session == nil || config == nil {
		return false
	}
	if ok := e.Session.Insert("GlobalConfig", config); ok {
		globalConfigMap[config.Name] = config
		globalConfigs = append(globalConfigs, config)
		return true
	}
	return false
}
