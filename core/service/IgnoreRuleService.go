package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type IgnoreRuleService struct {
	Session *db.MongoSession
}

//FindRuleByUser 获取规则
func (e *IgnoreRuleService) FindRuleByUser(userName string) (rules []*models.UserIgnoreRule) {
	coll := e.Session.GetCollection("IgnoreRule")
	if coll == nil {
		return nil
	}
	coll.Find(bson.M{"username": userName}).Select(nil).All(&rules)
	return
}

//FindRuleByMark 通过mark获取规则
func (e *IgnoreRuleService) FindRuleByMark(mark string) (rule *models.UserIgnoreRule) {
	coll := e.Session.GetCollection("IgnoreRule")
	if coll == nil {
		return nil
	}
	coll.Find(bson.M{"mark": mark}).Select(nil).One(rule)
	return
}

//AddRule 添加规则
func (e *IgnoreRuleService) AddRule(rule *models.UserIgnoreRule) {
	rule.Mark = rule.Labels.FastFingerprint().String()
	old := e.FindRuleByMark(rule.Mark)
	if old == nil {
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
