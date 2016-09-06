package user

import (
	"alertCenter/models"
	"fmt"
	"strings"
	"net/http"
	"bytes"
	"alertCenter/util"
	"github.com/astaxie/beego"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

const (
	apiVersion = "/api/v3"
	getUser = "/users?per_page=100"
	getGroup = "/groups?per_page=100"
)

var (
	gitlab = ""
	accessToken = ""
)

func init() {
	gitlab = beego.AppConfig.String("Gitlab")
	gitlab = strings.TrimSuffix(gitlab, "/")
	accessToken = beego.AppConfig.String("GitlabAccessToken")

}

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
	fmt.Println("In Gitlab Server SearchTeam")

	teams := []*models.Team{}

	page := 1
	for {
		url := gitlab + apiVersion + getGroup + "&page=" + strconv.Itoa(page)
		resp, err := GitlabApi("GET", url, nil)
		if err != nil {
			util.Error(err.Error())
			return nil, err
		}

		gitlabGroup := []GitlabGroup{}
		//fmt.Println("debug:", string(resp))

		err = json.Unmarshal(resp, &gitlabGroup)
		if err != nil {
			util.Error(err.Error())
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
	fmt.Println("In Gitlab Server SearchUsers")
	users := []*models.User{}

	page := 1
	for {
		url := gitlab + apiVersion + getUser + "&page=" + strconv.Itoa(page)
		resp, err := GitlabApi("GET", url, nil)
		if err != nil {
			util.Error(err.Error())
			return nil, err
		}

		gitlabusers := []GitlabUser{}

		//fmt.Println("debug:", string(resp))
		err = json.Unmarshal(resp, &gitlabusers)
		if err != nil {
			util.Error(err.Error())
			break
		}
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

func GitlabApi(method, url string, body []byte) ([]byte, error) {
	client := http.Client{}

	b := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		util.Error(err.Error())
		return []byte{}, err
	}

	req.Header.Set("PRIVATE-TOKEN", accessToken)
	//req.Header.Set("Authorization", "Bearer " + accessToken)
	resp, err := client.Do(req)
	if err != nil {
		util.Error(err.Error())
		return []byte{}, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error(err.Error())
		return []byte{}, err
	}
	defer resp.Body.Close()
	return respBody, nil
}


