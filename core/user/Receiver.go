package user

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

var (
	cacheTeams     map[string]*models.Team
	cacheApps      map[string]*models.APP
	cacheReceivers map[string]*models.Receiver
	cacheUsers     map[string]*models.User
	cacheTeamUsers map[string][]string
)

//Init 初始化用户关系缓存
func (r *Relation) Init() error {
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
		cacheTeamUsers = make(map[string][]string, 0)
	}

	ts := []*models.Team{}
	us := []*models.User{}
	source := beego.AppConfig.String("UserSource")
	sources := strings.Split(source, ",")

	var userServer UserInterface
	var err error
	for _, s := range sources {
		userServer, err = GetUserBySource(s)
		if err != nil {
			beego.Error(err.Error())
		}
		tmpTS, err := userServer.SearchTeams()
		if err != nil {
			return err
		}

		tmpUS, err := userServer.SearchUsers()
		if err != nil {
			return err
		}

		ts = append(ts, tmpTS...)
		us = append(us, tmpUS...)
	}
	beego.Info("load users number is " + strconv.Itoa(len(us)))
	beego.Info("load teams number is " + strconv.Itoa(len(ts)))
	if us != nil {
		for _, user := range us {
			if user.AvatarURL == "" {
				user.AvatarURL = "/static/img/default.jpg"
			}
			if user.Name != "" {
				cacheUsers[user.Name] = user
			}
		}
	}
	if ts != nil {
		for _, team := range ts {
			beego.Debug("load team " + team.Name)
			if team.Name != "" {
				cacheTeams[team.Name] = team
				users, err := userServer.GetUserByTeam(team.ID)
				if err != nil {
					return err
				}
				var completeUsers []string
				for _, user := range users {
					completeUsers = append(completeUsers, user.Name)
					beego.Debug("load user " + user.Name + " in team " + team.Name)
				}
				cacheTeamUsers[team.Name] = completeUsers
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

//GetAllAppInfo 获取全部app信息
func GetAllAppInfo() (apps []*models.APP, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", beego.AppConfig.String("cloudURI")+"/cloud-api/alert/apps", nil)
	if err != nil {
		beego.Error("create get appInfo request faild." + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	var AppInfo = make([]*models.APP, 0)
	err = json.Unmarshal(content, &AppInfo)
	if err != nil {
		beego.Error("Parse the appinfo data error ." + err.Error())
		return nil, err
	}
	return AppInfo, nil
}

//GetAppInfoByID 通过id查询app信息
func GetAppInfoByID(id string) (app *models.APP, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", beego.AppConfig.String("cloudURI")+"cloud-api/alert/app/"+id, nil)
	if err != nil {
		beego.Error("create get appInfo request faild." + err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("GET appinfo error," + err.Error())
		return nil, err
	}
	var AppInfo = models.APP{}
	err = json.Unmarshal(content, &AppInfo)
	if err != nil {
		beego.Error("Parse the appinfo data error ." + err.Error())
		return nil, err
	}
	return &AppInfo, nil
}

type Relation struct {
	session *db.MongoSession
}

func (r *Relation) SetTeam(team *models.Team) {
	team.ID = uuid.NewV4().String()
	// if team.WeTag == nil || team.WeTag.TagName == "" {
	// 	team.WeTag = GetWeTagByName(team.Name)
	// }
	cacheTeams[team.Name] = team
	r.session = db.GetMongoSession()
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

//GetUsersByTeam 根据组名获取用户
func (r *Relation) GetUsersByTeam(name string) (users []*models.User) {
	if cacheTeamUsers != nil {
		userNames := cacheTeamUsers[name]
		for _, userName := range userNames {
			user := cacheUsers[userName]
			if user != nil {
				users = append(users, user)
			}
		}
	}
	return
}

//GetUserByName 获取指定用户
func (r *Relation) GetUserByName(name string) *models.User {
	if name == "" {
		return nil
	}
	user := cacheUsers[name]
	if user != nil {
		return user
	} else {
		//暂时不实现
		return nil
	}
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
		var err error
		app, err = GetAppInfoByID(appID)
		if err != nil {
			return nil
		}
	}
	if app != nil {
		var us []string
		for _, mail := range app.Mails {
			user := FindUserByMail(mail)
			us = append(us, user.Name)
		}
		receiver = &models.Receiver{
			ID:        uuid.NewV4().String(),
			Name:      appID,
			UserNames: us,
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
		us := cacheTeamUsers[t.Name]
		receiver = &models.Receiver{
			ID:        uuid.NewV4().String(),
			Name:      team,
			UserNames: us,
		}
		cacheReceivers[team] = receiver
		return
	}
	return nil

}

//GetReceiver 获取receiver
func GetReceiver(labels models.Label) (receiver *models.Receiver) {

	if v, ok := labels.LabelSet["team"]; ok {
		return GetReceiverByTeam(string(v))
	}
	if v, ok := labels.LabelSet["container_label_app_id"]; ok {
		return GetReceiverByAPPID(string(v))
	}
	return nil
}
