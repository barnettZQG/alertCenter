package user

// 2016.9.20去除了receiver缓存，添加缓存自动刷新线程

import (
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

var (
	cacheTeams      map[string]*models.Team
	cacheApps       map[string]*models.APP
	cacheUsers      map[string]*models.User
	cacheTeamUsers  map[string][]string
	autoRefreshDura time.Duration
)

//Init 初始化用户关系缓存
func (r *Relation) Init() error {
	beego.Info("Relation init begin")
	defer beego.Info("Relation init over")
	err := initUser()
	if err != nil {
		return err
	}
	err = checkToken()
	if err != nil {
		return err
	}
	err = initApp()
	if err != nil {
		return err
	}

	autoRefreshTime, err := time.ParseDuration(beego.AppConfig.String("autoRefreshTime"))
	if err == nil {
		autoRefreshDura = autoRefreshTime
	} else {
		beego.Error(err)
		return err
	}
	go autoRefresh()
	return nil
}

func autoRefresh() {
	beego.Debug("open auto refresh users and teams work")
	var ticker = time.NewTicker(autoRefreshDura)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			beego.Debug("start auto refresh users and teams")
			initUser()
			checkToken()
			initApp()
			beego.Debug("end auto refresh users and teams")
		}
	}
}

//初始化用户及群组数据
func initUser() error {
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
			return err
		}
		tmpTS, err := userServer.SearchTeams()
		if err != nil {
			beego.Error(err.Error())
			return err
		}

		tmpUS, err := userServer.SearchUsers()
		if err != nil {
			beego.Error(err.Error())
			return err
		}

		ts = append(ts, tmpTS...)
		us = append(us, tmpUS...)
	}
	beego.Info("load users number is " + strconv.Itoa(len(us)))
	beego.Info("load teams number is " + strconv.Itoa(len(ts)))
	if us != nil {
		cacheUsersTmp := make(map[string]*models.User, 0)

		for _, user := range us {
			if user.AvatarURL == "" {
				user.AvatarURL = "/static/img/default.jpg"
			}
			if user.Name != "" {
				globalConfigService := service.GlobalConfigService{
					Session: db.GetMongoSession(),
				}
				if globalConfigService.Session != nil {
					defer globalConfigService.Session.Close()
				}
				if ok, _ := globalConfigService.CheckExist("IsAdmin", user.Name); ok {
					user.IsAdmin = true
				}
				cacheUsersTmp[user.Name] = user
			}
		}
		cacheUsers = cacheUsersTmp
	}
	if ts != nil {

		cacheTeamsTmp := make(map[string]*models.Team, 0)
		cacheTeamUsersTmp := make(map[string][]string, 0)
		for _, team := range ts {
			//beego.Debug("load team " + team.Name)
			if team.Name != "" {
				cacheTeamsTmp[team.Name] = team
				users, err := userServer.GetUserByTeam(team.ID)
				if err != nil {
					beego.Error(err.Error())
					return err
				}
				var completeUsers []string
				for _, user := range users {
					completeUsers = append(completeUsers, user.Name)
					beego.Debug("load user " + user.Name + " in team " + team.Name)
				}
				cacheTeamUsersTmp[team.Name] = completeUsers
			}
		}
		cacheTeams = cacheTeamsTmp
		cacheTeamUsers = cacheTeamUsersTmp
	}
	return nil
}

//初始化app数据
func initApp() error {
	if beego.AppConfig.String("cloudURI") == "" {
		return nil
	}
	as, err := GetAllAppInfo()
	if err != nil {
		return err
	}
	beego.Info("load apps number is " + strconv.Itoa(len(as)))

	cacheAppsTmp := make(map[string]*models.APP, 0)

	for _, app := range as {
		if app.ID != "" {
			cacheAppsTmp[app.ID] = app
		}
	}
	cacheApps = cacheAppsTmp
	return nil
}

//检查用户默认token是否存在
func checkToken() error {
	service := &service.TokenService{
		Session: db.GetMongoSession(),
	}
	if service.Session != nil {
		defer service.Session.Close()
	}
	for _, user := range cacheUsers {
		token := service.GetDefaultToken(user.Name)
		if token == nil {
			service.CreateToken("default", user.Name)
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

//RefreshCache 更新缓存
func (r *Relation) RefreshCache() {
	beego.Debug("start refresh users and teams")
	initUser()
	checkToken()
	initApp()
	beego.Debug("end  refresh users and teams")
}

//SetTeam 老方法添加team
func (r *Relation) SetTeam(team *models.Team) {
	team.ID = uuid.NewV4().String()
	// if team.WeTag == nil || team.WeTag.TagName == "" {
	// 	team.WeTag = GetWeTagByName(team.Name)
	// }
	cacheTeams[team.Name] = team
	r.session = db.GetMongoSession()
	if r.session != nil {
		defer r.session.Close()
	}
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
		// return &models.User{
		// 	Name: name,
		// }
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

	}
	return
}

//GetReceiverByTeam 通过team name 获取receiver
func GetReceiverByTeam(team string) (receiver *models.Receiver) {
	t := cacheTeams[team]
	if t != nil {
		us := cacheTeamUsers[t.Name]
		receiver = &models.Receiver{
			ID:        uuid.NewV4().String(),
			Name:      team,
			UserNames: us,
		}
		return
	}
	return nil

}

//GetReceiverByUser 通过user name 获取receiver
func GetReceiverByUser(user string) (receiver *models.Receiver) {

	users := strings.Split(user, ",")
	var us []string
	for _, u := range users {
		t := cacheUsers[u]
		if t != nil {
			us = append(us, t.Name)
		}
	}
	receiver = &models.Receiver{
		ID:        uuid.NewV4().String(),
		Name:      "user_" + user,
		UserNames: us,
	}
	return
}

//GetReceiver 获取receiver
func GetReceiver(labels models.Label) (receiver *models.Receiver) {

	if v, ok := labels.LabelSet["user"]; ok {
		return GetReceiverByUser(string(v))
	}
	if v, ok := labels.LabelSet["container_label_app_id"]; ok {
		return GetReceiverByAPPID(string(v))
	}
	if v, ok := labels.LabelSet["team"]; ok {
		return GetReceiverByTeam(string(v))
	}
	return nil
}
