package controllers

import (
	"alertCenter/controllers/session"
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"alertCenter/util"
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
)

type IgnoreRuleAPIControll struct {
	APIBaseController
}

//AddRule 添加自定义忽略规则
func (e *IgnoreRuleAPIControll) AddRule() {
	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		beego.Debug("ignoreRule:" + string(data))
		var rule *models.UserIgnoreRule = &models.UserIgnoreRule{}
		err := json.Unmarshal(data, rule)
		if err == nil {
			session := db.GetMongoSession()
			if session != nil {
				defer session.Close()
			}

			ruleService := &service.IgnoreRuleService{
				Session: session,
			}
			ruleService.AddRule(rule)
			e.Data["json"] = util.GetSuccessReJson(rule)
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

//GetRulesByUser 获取用户的规则
func (e *IgnoreRuleAPIControll) GetRulesByUser() {
	user := e.Ctx.Input.Header("user")
	if user == "" {
		e.Data["json"] = util.GetErrorJson("please certification")
		e.ServeJSON()
	} else {
		session := db.GetMongoSession()
		if session != nil {
			defer session.Close()
		}
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

//AddRuleByAlert 添加某alert的忽略规则
func (e *IgnoreRuleAPIControll) AddRuleByAlert() {
	ID := e.GetString(":mark")
	sess, err := session.GetSession(e.Ctx)
	if err != nil {
		beego.Error("get session error:", err)
		e.Data["json"] = util.GetErrorJson("please certification")
		e.ServeJSON()
	}
	user := sess.Get(session.SESSION_USER)
	if ID == "" {
		e.Data["json"] = util.GetErrorJson("api error,mark is not provided")
	} else {
		session := db.GetMongoSession()
		if session != nil {
			defer session.Close()
		}
		alertService := &service.AlertService{
			Session: session,
		}
		alert := alertService.FindByID(ID)
		if alert == nil {
			e.Data["json"] = util.GetErrorJson("alertID is not exit")
		} else {
			rule := &models.UserIgnoreRule{
				Labels:   alert.Labels,
				StartsAt: time.Now(),
				UserName: user.(*models.User).Name,
			}
			ruleService := &service.IgnoreRuleService{
				Session: session,
			}
			ruleService.AddRule(rule)
			e.Data["json"] = util.GetSuccessJson("add user ignore rule success")
		}
	}
	e.ServeJSON()
}

//DeleteRule 删除rule
func (e *IgnoreRuleAPIControll) DeleteRule() {
	user := e.Ctx.Input.Header("user")
	ruleID := e.GetString(":ruleID")
	beego.Debug("delete ignoreRule,user:" + user)
	if user == "" {
		e.Data["json"] = util.GetErrorJson("please certification")
		e.ServeJSON()
	} else {
		service := &service.IgnoreRuleService{
			Session: db.GetMongoSession(),
		}
		if service.Session != nil {
			defer service.Session.Close()
		}

		if ok := service.DeleteRule(ruleID, user); ok {
			e.Data["json"] = util.GetSuccessJson("Delete rule success")
		} else {
			e.Data["json"] = util.GetFailJson("Delete rule faild")
		}
		e.ServeJSON()
	}
}
