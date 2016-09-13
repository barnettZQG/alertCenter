package user

import (
	"alertCenter/models"
	"encoding/json"
	"strconv"
	"strings"

	"alertCenter/core/gitlab"
	"github.com/astaxie/beego"
	"fmt"
)

type GitlabUser struct {
	Id       int
	Username string
	Email    string
	Name     string
	State    string
	Is_admin bool
	Bio      string
}

type GitlabGroup struct {
	Id          int
	Name        string
	Path        string
	Description string
}

type GitlabServer struct {
}

func (e *GitlabServer) SearchTeams() ([]*models.Team, error) {
	beego.Info("In Gitlab Server SearchTeam")

	teams := []*models.Team{}

	page := 1
	for {
		url := gitlab.GetGitlabUrl() + gitlab.ApiVersion + gitlab.GetGroup + "?per_page=100&page=" + strconv.Itoa(page)
		resp, err := gitlab.GitlabApi("GET", url, nil)
		if err != nil {
			beego.Error(err.Error())
			return nil, err
		}

		gitlabGroup := []GitlabGroup{}
		//fmt.Println("debug:", string(resp))

		err = json.Unmarshal(resp, &gitlabGroup)
		if err != nil {
			beego.Error(err.Error())
			break
		}
		if len(gitlabGroup) == 0 {
			break
		}

		for _, g := range gitlabGroup {
			tmp := &models.Team{}
			tmp.ID = strconv.Itoa(g.Id)
			tmp.Name = g.Name
			teams = append(teams, tmp)
		}
		page = page + 1
	}

	return teams, nil
}

func (e *GitlabServer) SearchUsers() ([]*models.User, error) {
	beego.Info("In Gitlab Server SearchUsers")
	users := []*models.User{}

	page := 1
	for {
		url := gitlab.GetGitlabUrl() + gitlab.ApiVersion + gitlab.GetUser + "?per_page=100&page=" + strconv.Itoa(page)

		resp, err := gitlab.GitlabApi("GET", url, nil)
		if err != nil {
			beego.Error(err.Error())
			return nil, err
		}

		gitlabusers := []GitlabUser{}

		//fmt.Println("debug:", string(resp))
		err = json.Unmarshal(resp, &gitlabusers)
		if err != nil {
			fmt.Println(string(resp))
			beego.Error(err.Error())
			break
		}
		beego.Debug("SearchUsers, this loop get user:", len(gitlabusers))
		if len(gitlabusers) == 0 {
			break
		}

		for _, u := range gitlabusers {
			if u.State != "active" {
				continue
			}
			tmp := &models.User{}
			tmp.ID = strconv.Itoa(u.Id)
			tmp.Name = u.Username
			tmp.Mail = u.Email
			users = append(users, tmp)
		}
		page = page + 1

	}

	return users, nil
}

func (e *GitlabServer) GetUserByTeam(id string) ([]*models.User, error) {
	beego.Info("In Gitlab Server GetUserByTeam")
	users := []*models.User{}

	page := 1
	for {

		guUrl := strings.Replace(gitlab.GetGroupUsers, ":id", id, -1)

		url := gitlab.GetGitlabUrl() + gitlab.ApiVersion + guUrl + "?per_page=100&page=" + strconv.Itoa(page)
		resp, err := gitlab.GitlabApi("GET", url, nil)
		if err != nil {
			beego.Error(err.Error())
			return nil, err
		}

		gitlabusers := []GitlabUser{}

		//fmt.Println("debug:", string(resp))
		err = json.Unmarshal(resp, &gitlabusers)
		if err != nil {
			beego.Error(err.Error())
			break
		}

		for _, u := range gitlabusers {
			if u.State != "active" {
				continue
			}
			tmp := &models.User{}
			tmp.ID = strconv.Itoa(u.Id)
			tmp.Name = u.Username
			users = append(users, tmp)
		}
		page = page + 1
		//此处不严谨，根据api版本可能不成立
		if len(gitlabusers) < 100 {
			break
		}
	}

	return users, nil
}



