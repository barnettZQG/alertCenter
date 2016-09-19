package controllers

import (
	"alertCenter/core/db"
	"alertCenter/core/user"
	"alertCenter/models"
	"alertCenter/util"
	"encoding/json"

	"github.com/astaxie/beego"
)

type TeamAPIController struct {
	APIBaseController
	session *db.MongoSession
}

func (e *TeamAPIController) GetTeams() {
	e.session = db.GetMongoSession()
	if e.session != nil {
		defer e.session.Close()
	}
	if e.session == nil {
		e.Data["json"] = util.GetFailJson("get database session faild.")
	} else {
		relation := &user.Relation{}
		teams := relation.GetAllTeam()
		if teams == nil {
			e.Data["json"] = util.GetFailJson("There is no info of team")
		} else {
			e.Data["json"] = util.GetSuccessReJson(teams)
		}
	}
	e.ServeJSON()
}

func (e *TeamAPIController) AddTeam() {
	data := e.Ctx.Input.RequestBody
	if data != nil && len(data) > 0 {
		var team *models.Team = &models.Team{}
		err := json.Unmarshal(data, team)
		if err == nil {
			relation := &user.Relation{}
			relation.SetTeam(team)
			e.Data["json"] = util.GetSuccessJson("receive team info success")
		} else {
			beego.Error("Parse the received message to teams faild." + err.Error())
			e.Data["json"] = util.GetFailJson("Parse the received message to teams faild.")
		}
	} else {
		beego.Error("receive a unknow data")
		e.Data["jaon"] = util.GetErrorJson("receive a unknow data")
	}
	e.ServeJSON()
}

//type TeamUser struct {
//	Team *models.Team
//	User []*models.User
//}
//
//func (e *TeamAPIController) GetTeamUsers() {
//	e.session = db.GetMongoSession()
//	if e.session == nil {
//		e.Data["json"] = util.GetFailJson("get database session faild.")
//	} else {
//		relation := &user.Relation{}
//		teams := relation.GetAllTeam()
//		if teams == nil {
//			e.Data["json"] = util.GetFailJson("There is no info of team")
//		} else {
//			var result []*TeamUser
//			for _, team := range teams {
//				users := relation.GetUsersByTeam(team.Name)
//				teamUser := &TeamUser{
//					Team: team,
//					User: users,
//				}
//				result = append(result, teamUser)
//			}
//			e.Data["json"] = util.GetSuccessReJson(result)
//		}
//	}
//	e.ServeJSON()
//}
