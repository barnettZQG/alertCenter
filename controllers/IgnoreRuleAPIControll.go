package controllers

import (
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"alertCenter/util"
	"encoding/json"

	"github.com/astaxie/beego"
)

type IgnoreRuleAPIControll struct {
	beego.Controller
}

func (e *IgnoreRuleAPIControll) AddRule() {
	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var rule *models.UserIgnoreRule = &models.UserIgnoreRule{}
		err := json.Unmarshal(data, rule)
		if err == nil {
			session := db.GetMongoSession()
			defer session.Close()
			ruleService := &service.IgnoreRuleService{
				Session: session,
			}
			ruleService.AddRule(rule)
			e.Data["json"] = util.GetSuccessJson("receive user ignore rule info success")
		} else {
			beego.Error("Parse the received user ignore rule faild." + err.Error())
			e.Data["json"] = util.GetFailJson("Parse the received user ignore rule faild.")
		}
	} else {
		beego.Error("receive a unknow data")
		e.Data["jaon"] = util.GetErrorJson("receive a unknow data")
	}
	e.ServeJSON()
}

func (e *IgnoreRuleAPIControll) GetRulesByUser() {
	user := e.GetString(":user")
	if user == "" {
		e.Data["json"] = util.GetErrorJson("api error,userID is not provided")
	} else {
		session := db.GetMongoSession()
		defer session.Close()
		ruleService := &service.IgnoreRuleService{
			Session: session,
		}
		rules := ruleService.FindRuleByUser(user)
		if rules != nil {
			e.Data["json"] = util.GetSuccessReJson(rules)
		} else {
			e.Data["json"] = util.GetFailJson("userID is not exit or this user have not rules")
		}
	}
	e.ServeJSON()
}
