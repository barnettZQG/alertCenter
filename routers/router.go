package routers

import (
	"alertCenter/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/logout", &controllers.MainController{}, "post:Logout")
	beego.Router("/transit", &controllers.MainController{}, "get:Transit")
	beego.Router("/alerts", &controllers.AlertController{}, "get:AlertList")
	beego.Router("/alertsCurrent", &controllers.AlertController{}, "get:AlertsCurrent")
	beego.Router("/teams", &controllers.TeamController{}, "get:GetTeams")
	beego.Router("teamUsers", &controllers.TeamController{}, "get:GetTeamUsers")
	beego.Router("/user/:userName", &controllers.UserController{}, "get:UserHome")
	beego.Router("/token/addToken", &controllers.TokenController{}, "post:AddToken")

	beego.Router("/api/teams", &controllers.TeamAPIController{}, "get:GetTeams")
	beego.Router("/api/addTeam", &controllers.TeamAPIController{}, "post:AddTeam")
	beego.Router("/api/receive", &controllers.APIController{}, "post:Receive")
	beego.Router("/api/v1/alerts", &controllers.PrometheusAPI{}, "post:ReceivePrometheus")
	//beego.Router("/api/getTag", &controllers.APIController{}, "get:AddTag")
	beego.Router("/api/alert/handle/:ID/:type", &controllers.APIController{}, "post:HandleAlert")
	beego.Router("/api/ignoreRules", &controllers.IgnoreRuleAPIControll{}, "get:GetRulesByUser")
	beego.Router("/api/ignoreRule/:ruleID", &controllers.IgnoreRuleAPIControll{}, "delete:DeleteRule")
	beego.Router("/api/alerts", &controllers.APIController{}, "get:GetAlerts")
	beego.Router("/api/addIgnoreRule", &controllers.IgnoreRuleAPIControll{}, "post:AddRule")
	beego.Router("/api/ignoreAlert/:mark", &controllers.IgnoreRuleAPIControll{}, "post:AddRuleByAlert")
	beego.Router("/api/projects", &controllers.TokenAPIController{}, "get:GetAllToken")
	beego.Router("/api/project/:project", &controllers.TokenAPIController{}, "delete:DeleteToken")
	//外部通知开关控制
	beego.Router("/api/noticeOn/:status", &controllers.APIController{}, "post:SetNoticeMode")
	beego.Router("/api/noticeOn", &controllers.APIController{}, "get:GetNoticeMode")
	//白名单ip控制
	beego.Router("/api/trustIP", &controllers.APIController{}, "post:AddTrustIP")
	beego.Router("/api/trustIP", &controllers.APIController{}, "get:GetTrustIP")
	beego.Router("/api/trustIP/:ID", &controllers.APIController{}, "delete:DeleteTrustIP")

	beego.Router("/api/refreshCache", &controllers.APIController{}, "post:RefreshCache")
}
