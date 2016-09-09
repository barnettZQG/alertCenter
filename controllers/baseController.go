package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	"net/http"
	//"fmt"
)

type BaseController struct {
	beego.Controller
}

var globalSessions *session.Manager

func init() {
	mangeConfig := &session.ManagerConfig{CookieLifeTime:3600, CookieName:"gosessionid", Gclifetime:3600,EnableSetCookie:true}

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
	if sessCode == nil && paramCode == "" {
		//fmt.Printf("session %#v\n",sess)
		redirct := beego.AppConfig.String("Gitlab") + "/oauth/authorize?response_type=code&client_id=1b44e5046e51199065f89f0ad2a59385978cd025bedb1353c9148f7496907804&redirect_uri=http%3A%2F%2Flocalhost:8888"
		http.Redirect(this.Ctx.ResponseWriter, this.Ctx.Request, redirct, http.StatusTemporaryRedirect)
		return
	}  else {
		//fmt.Printf("session %#v\n",sess)
		err := sess.Set("code",paramCode)
		if err!=nil {
			beego.Error(err)
		}
		//beego.Debug("Have code.", "sessCode",sessCode,"paramCode",paramCode)
	}
}