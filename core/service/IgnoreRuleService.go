package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"time"

	"gopkg.in/mgo.v2/bson"

	uuid "github.com/satori/go.uuid"
)

type IgnoreRuleService struct {
	Session *db.MongoSession
}

//FindRuleByUser 获取规则
func (e *IgnoreRuleService) FindRuleByUser(userID string) (rules []*models.UserIgnoreRule) {
	coll := e.Session.GetCollection("IgnoreRule")
	if coll == nil {
		return nil
	}
	coll.Find(bson.M{"userid": userID}).Select(nil).All(&rules)
	return
}

//AddRule 添加规则
func (e *IgnoreRuleService) AddRule(rule *models.UserIgnoreRule) {
	rule.AddTime = time.Now()
	rule.RuleID = uuid.NewV4().String()
	rule.IsLive = true
	e.Session.Insert("IgnoreRule", rule)
}
