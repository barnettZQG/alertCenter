package controllers

import (
	"alertCenter/models"
	"alertCenter/core/gitlab"
	"alertCenter/util"
	"fmt"
)

type TeamController struct {
	BaseController
}

func (e *TeamController) GetTeams() {
	e.TplName = "teams.html"
}

type TeamUser struct {
	Team *models.Team
	User []*models.User
}

func (e *TeamController) GetTeamUsers() {
	groups, err := gitlab.GetGroupsByUsername(e.BaseController.Username)
	if err != nil {

	}
	var result []*TeamUser
	for _, group := range groups {
		team := gitlab.ConvertGitlabGroupToAlertModel(group)
		users, _ := gitlab.GetUsersByTeam(e.BaseController.Username, team.ID)
		us := gitlab.ConvertGitlabUsers(users)
		teamUser := &TeamUser{
			Team: team,
			User: us,
		}
		result = append(result, teamUser)
	}

	fmt.Println("in GetTeamUsers and result;",result)

	e.Data["json"] = util.GetSuccessReJson(result)
	e.ServeJSON()
}
