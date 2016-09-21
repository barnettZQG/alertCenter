package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type IgnoreRuleService struct {
	Session *db.MongoSession
}

//FindRuleByUser 获取规则
func (e *IgnoreRuleService) FindRuleByUser(userName string) (rules []*models.UserIgnoreRule, err error) {
	coll := e.Session.GetCollection("IgnoreRule")
	if coll == nil {
		return nil, fmt.Errorf("get IgnoreRule collection error")
	}
	err = coll.Find(bson.M{"username": userName}).Select(nil).All(&rules)
	return
}

//FindRuleByMark 通过mark获取规则
func (e *IgnoreRuleService) FindRuleByMark(mark string) (rule *models.UserIgnoreRule, err error) {
	coll := e.Session.GetCollection("IgnoreRule")
	if coll == nil {
		return nil, fmt.Errorf("get IgnoreRule collection error")
	}
	err = coll.Find(bson.M{"mark": mark}).Select(nil).One(rule)
	return
}

//AddRule 添加规则
func (e *IgnoreRuleService) AddRule(rule *models.UserIgnoreRule) {
	rule.Mark = rule.Labels.FastFingerprint().String()
	_, err := e.FindRuleByMark(rule.Mark)
	if err != nil && err.Error() == mgo.ErrNotFound.Error() {
		rule.AddTime = time.Now()
		rule.RuleID = uuid.NewV4().String()
		rule.IsLive = true
		e.Session.Insert("IgnoreRule", rule)
	}
}

//DeleteRule 删除规则
func (e *IgnoreRuleService) DeleteRule(ruleID string, user string) bool {
	col := e.Session.GetCollection("IgnoreRule")
	if col == nil {
		beego.Error("get collection IgnoreRule error ")
		return false
	}
	err := col.Remove(bson.M{"ruleid": ruleID, "username": user})
	if err != nil {
		beego.Error("delete IgnoreRule error ,", err.Error())
		return false
	}
	return true
}
