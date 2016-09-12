package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	//"fmt"
	"alertCenter/core/gitlab"
	"fmt"
)

type BaseController struct {
	beego.Controller
}

var globalSessions *session.Manager

const (
	SESSION_USER     = "user"
	SESSION_USERNAME = "username"
)

func init() {

	mangeConfig := &session.ManagerConfig{CookieLifeTime: 3600, CookieName: "gosessionid", Gclifetime: 3600, EnableSetCookie: true}

	globalSessions, _ = session.NewManager("memory", mangeConfig)
	go globalSessions.GC()
}

func (this *BaseController) Prepare() {
	sess, err := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	if err != nil {
		beego.Error("get session error:", err)
		return
	}

	sessUsername := sess.Get(SESSION_USERNAME)
	paramCode := this.GetString("code")
	fmt.Printf("paramCode: %s,session: %#v\n", paramCode, sess)

	if sessUsername == nil && paramCode == "" {
		fmt.Println("in sessUsername == nil && paramCode == nil")
		redirct := gitlab.GetGitlabOAuthUrl()
		http.Redirect(this.Ctx.ResponseWriter, this.Ctx.Request, redirct, http.StatusTemporaryRedirect)
		return
	} else if sessUsername == nil && paramCode != "" {
		fmt.Println("in sessUsername == nil && paramCode != nil")
		access, err := gitlab.GetGitlabAccessToken(paramCode)
		if err != nil {
			beego.Error(err)
			return
		}
		user, err := gitlab.GetCurrentUserWithToken(access.AccessToken)
		if err != nil {
			beego.Error(err)
			return
		}
		err = sess.Set(SESSION_USER, user)
		if err != nil {
			beego.Error(err)
			return
		}
		// check if the code is right.
		err = sess.Set(SESSION_USERNAME, user.Username)
		if err != nil {
			beego.Error(err)
			return
		}
		gitlab.Tokens.Add(user.Username, access)
	} else {
		fmt.Println("in sessUsername != nil && paramCode != nil")
		u := sess.Get(SESSION_USER)
		n := sess.Get(SESSION_USERNAME)
		fmt.Println("username:", n, "user:", u)
		// Already login
		//beego.Debug("Have code.", "sessCode",sessCode,"paramCode",paramCode)
	}
}
