package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"alertCenter/core/interfaces"
	"alertCenter/models"
	"alertCenter/util"

	"github.com/astaxie/beego"
	"github.com/prometheus/common/model"
	uuid "github.com/satori/go.uuid"
)

var (
	cacheTeams     map[string]*models.Team
	cacheApps      map[string]*models.APP
	cacheReceivers map[string]*models.Receiver
	cacheUsers     map[string]*models.User
	cacheTeamUsers map[string][]*models.User
)

func (r *Relation) Init(usermanager interfaces.UserManager) error {
	beego.Info("Relation init begin")
	defer beego.Info("Relation init over")
	if cacheTeams == nil {
		cacheTeams = make(map[string]*models.Team, 0)
	}
	if cacheApps == nil {
		cacheApps = make(map[string]*models.APP, 0)
	}
	if cacheReceivers == nil {
		cacheReceivers = make(map[string]*models.Receiver, 0)
	}
	if cacheUsers == nil {
		cacheUsers = make(map[string]*models.User, 0)
	}
	if cacheTeamUsers == nil {
		cacheTeamUsers = make(map[string][]*models.User, 0)
	}
	ts, err := usermanager.SearchTeams()
	if err != nil {
		return err
	}
	us, err := usermanager.SearchUsers()
	if err != nil {
		return err
	}
	beego.Info("load users number is " + strconv.Itoa(len(us)))
	beego.Info("load teams number is " + strconv.Itoa(len(ts)))
	if ts != nil {
		for _, team := range ts {
			if team.Name != "" {
				cacheTeams[team.Name] = team
			}
		}
	}
	if us != nil {
		for _, user := range us {
			if user.Name != "" {
				cacheUsers[user.Name] = user
				if us := cacheTeamUsers[user.TeamID]; us != nil {
					us = append(us, user)
					cacheTeamUsers[user.TeamID] = us
				} else {
					var us []*models.User
					us = append(us, user)
					cacheTeamUsers[user.TeamID] = us
				}
			}
		}
	}
	as, err := GetAllAppInfo()
	if err != nil {
		return err
	}
	beego.Info("load apps number is " + strconv.Itoa(len(as)))
	for _, app := range as {
		if app.ID != "" {
			cacheApps[app.ID] = app
		}
	}
	return nil
}

func GetAllAppInfo() (apps []*models.APP, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", beego.AppConfig.String("cloudURI")+"/cloud-api/alert/apps", nil)
	if err != nil {
		util.Error("create get appInfo request faild." + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		util.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	var AppInfo = make([]*models.APP, 0)
	err = json.Unmarshal(content, &AppInfo)
	if err != nil {
		util.Error("Parse the appinfo data error ." + err.Error())
		return nil, err
	}
	return AppInfo, nil
}
func GetAppInfoById(id string) (app *models.APP, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", beego.AppConfig.String("cloudURI")+"cloud-api/alert/app/"+id, nil)
	if err != nil {
		util.Error("create get appInfo request faild." + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		util.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	var AppInfo = models.APP{}
	err = json.Unmarshal(content, &AppInfo)
	if err != nil {
		util.Error("Parse the appinfo data error ." + err.Error())
		return nil, err
	}
	return &AppInfo, nil
}

type Relation struct {
	session *MongoSession
}

func (r *Relation) SetTeam(team *models.Team) {
	team.ID = uuid.NewV4().String()
	// if team.WeTag == nil || team.WeTag.TagName == "" {
	// 	team.WeTag = GetWeTagByName(team.Name)
	// }
	cacheTeams[team.Name] = team
	r.session = GetMongoSession()
	defer r.session.Close()
	r.session.Insert("team", team)
}

//GetAllTeam 获取全部团队
func (r *Relation) GetAllTeam() (result []*models.Team) {
	if cacheTeams != nil {
		for _, v := range cacheTeams {
			result = append(result, v)
		}
	}
	return
}

//GetAllUser 获取全部用户
func (r *Relation) GetAllUser() (result []*models.User) {
	if cacheUsers != nil {
		for _, v := range cacheUsers {
			result = append(result, v)
		}
	}
	return
}

//FindUserByMail 通过邮箱获取用户
func FindUserByMail(mail string) *models.User {
	if mail == "" {
		return nil
	}
	if cacheUsers != nil {
		for _, v := range cacheUsers {
			if v.Mail == mail {
				return v
			}
		}
	}
	return nil
}

//GetReceiverByAPPID 通过appid获取receiver
func GetReceiverByAPPID(appID string) (receiver *models.Receiver) {
	if appID == "" {
		return
	}
	receiver = cacheReceivers[appID]
	if receiver != nil {
		return
	}
	app := cacheApps[appID]
	if app == nil {
		app, err = GetAppInfoById(appID)
		if err != nil {
			return nil
		}
	}
	if app != nil {
		var us = make([]*models.User, 0)
		for _, mail := range app.Mails {
			user := FindUserByMail(mail)
			us = append(us, user)
		}
		receiver = &models.Receiver{
			ID:    uuid.NewV4().String(),
			Name:  appID,
			Users: us,
		}
		cacheReceivers[appID] = receiver
	}
	return
}

//GetReceiverByTeam 通过team name 获取receiver
func GetReceiverByTeam(team string) (receiver *models.Receiver) {
	re := cacheReceivers[team]
	if re != nil {
		return re
	}
	t := cacheTeams[team]
	if t != nil {
		us := cacheTeamUsers[t.ID]
		receiver = &models.Receiver{
			ID:    uuid.NewV4().String(),
			Name:  team,
			Users: us,
		}
		cacheReceivers[team] = receiver
		return
	}
	return nil

}

//GetReceiver 获取receiver
func GetReceiver(labels model.LabelSet) (receiver *models.Receiver) {

	if v, ok := labels["team"]; ok {
		return GetReceiverByTeam(string(v))
	}
	if v, ok := labels["container_label_app_id"]; ok {
		return GetReceiverByAPPID(string(v))
	}
	return nil
}
