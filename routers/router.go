package routers

import (
	"alertCenter/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/transit", &controllers.MainController{}, "get:Transit")
	beego.Router("/alerts", &controllers.AlertController{}, "get:AlertList")
	beego.Router("/teams", &controllers.TeamController{}, "get:GetTeams")
	beego.Router("/user/:userName", &controllers.UserController{}, "get:UserHome")
	beego.Router("/token/addToken", &controllers.TokenController{}, "post:AddToken")
	beego.Router("/api/teams", &controllers.TeamAPIController{}, "get:GetTeams")
	beego.Router("/api/teamUsers", &controllers.TeamAPIController{}, "get:GetTeamUsers")
	beego.Router("/api/addTeam", &controllers.TeamAPIController{}, "post:AddTeam")
	beego.Router("/api/receive", &controllers.APIController{}, "post:Receive")
	beego.Router("/api/v1/alerts", &controllers.APIController{}, "post:Receive")
	//beego.Router("/api/getTag", &controllers.APIController{}, "get:AddTag")
	beego.Router("/api/alert/handle/:ID/:type", &controllers.APIController{}, "post:HandleAlert")
	beego.Router("/api/ignoreRules", &controllers.IgnoreRuleAPIControll{}, "get:GetRulesByUser")
	beego.Router("/api/alerts", &controllers.APIController{}, "get:GetAlerts")
	beego.Router("/api/addIgnoreRule", &controllers.IgnoreRuleAPIControll{}, "post:AddRule")
	beego.Router("/api/ignoreAlert/:mark", &controllers.IgnoreRuleAPIControll{}, "post:AddRuleByAlert")
}
