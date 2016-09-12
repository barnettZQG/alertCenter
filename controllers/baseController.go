package controllers

import (
	"github.com/astaxie/beego"

	"net/http"
	//"fmt"
	"fmt"
	"alertCenter/core/gitlab"
	"alertCenter/controllers/session"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {

	sess, err := session.GetSession(this.Ctx)
	if err != nil {
		beego.Error(err)
		return
	}

	sessUsername := sess.Get(session.SESSION_USERNAME)
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
		err = sess.Set(session.SESSION_USER, user)
		if err != nil {
			beego.Error(err)
			return
		}
		// check if the code is right.
		err = sess.Set(session.SESSION_USERNAME, user.Username)
		if err != nil {
			beego.Error(err)
			return
		}
		gitlab.Tokens.Add(user.Username, access)
	} else {
		fmt.Println("in sessUsername != nil && paramCode != nil")
		u := sess.Get(session.SESSION_USER)
		n := sess.Get(session.SESSION_USERNAME)
		fmt.Println("username:", n, "user:", u)
		// Already login
		//beego.Debug("Have code.", "sessCode",sessCode,"paramCode",paramCode)
	}
}

