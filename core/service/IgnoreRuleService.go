package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
)

type IgnoreRuleService struct {
	Session *db.MongoSession
}

//FindRuleByUser 获取规则
func (e *IgnoreRuleService) FindRuleByUser() (rule *models.UserIgnoreRule) {
	return
}
