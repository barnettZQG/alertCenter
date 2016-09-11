package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	"net/http"
	//"fmt"
	"fmt"
)

type BaseController struct {
	beego.Controller
}

var clientId = "1b44e5046e51199065f89f0ad2a59385978cd025bedb1353c9148f7496907804"
var sercet = "15b139748d066635620491ebfcd0349b5a4b5457ba9d7728e2e3105621b227cc"
var globalSessions *session.Manager

func init() {
	mangeConfig := &session.ManagerConfig{CookieLifeTime:3600, CookieName:"gosessionid", Gclifetime:3600, EnableSetCookie:true}

	globalSessions, _ = session.NewManager("memory", mangeConfig)
	go globalSessions.GC()
}

func (this *BaseController) Prepare() {
	sess, err := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	if err != nil {
		beego.Error("get session error:", err)
		return
	}

	sessCode := sess.Get("code")
	paramCode := this.GetString("code")
	fmt.Printf("paramCode: %s,session: %#v\n",paramCode, sess)

	if sessCode == nil && paramCode == "" {
		redirct := beego.AppConfig.String("Gitlab") + "/oauth/authorize?response_type=code&client_id=1b44e5046e51199065f89f0ad2a59385978cd025bedb1353c9148f7496907804&redirect_uri=http%3A%2F%2Flocalhost:8888"
		http.Redirect(this.Ctx.ResponseWriter, this.Ctx.Request, redirct, http.StatusTemporaryRedirect)
		return
	}else if sessCode == nil && paramCode != ""{
		// check if the code is right.
		err := sess.Set("code", paramCode)
		if err != nil {
			beego.Error(err)
		}
	}else {
		// Already login
		//beego.Debug("Have code.", "sessCode",sessCode,"paramCode",paramCode)
	}
}

func GetGitlabAccessToken(code string) (error) {
	url := beego.AppConfig.String("Gitlab") + "/oauth/token?client_id=" + clientId + "&client_secret=" + sercet + "&code=" + code + "&grant_type=authorization_code&redirect_uri=http://localhost:8888"
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	_ = resp
	return nil
}