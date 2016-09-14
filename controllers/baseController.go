package controllers

import (
	"alertCenter/controllers/session"
	"alertCenter/core/db"
	"alertCenter/core/gitlab"
	"alertCenter/core/service"
	"alertCenter/core/user"
	"net/http"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	Username string
}

func (this *BaseController) Prepare() {

	sess, err := session.GetSession(this.Ctx)
	if err != nil {
		beego.Error(err)
		return
	}
	defer sess.SessionRelease(this.Ctx.ResponseWriter)

	sessUsername := sess.Get(session.SESSION_USERNAME)
	paramCode := this.GetString("code")
	//fmt.Printf("paramCode: %s,session: %#v\n", paramCode, sess)

	if sessUsername == nil && paramCode == "" {

		sess.Set("redirect", this.Ctx.Request.URL.String()) //为了再次访问的重定向
		//fmt.Println("in sessUsername == nil && paramCode == nil")
		redirct := gitlab.GetGitlabOAuthUrl()
		http.Redirect(this.Ctx.ResponseWriter, this.Ctx.Request, redirct, http.StatusTemporaryRedirect)
		return
	} else if sessUsername == nil && paramCode != "" {
		//fmt.Println("in sessUsername == nil && paramCode != nil")
		access, err := gitlab.GetGitlabAccessToken(paramCode)
		if err != nil {
			beego.Error(err)
			return
		}
		u, err := gitlab.GetCurrentUserWithToken(access.AccessToken)
		if err != nil {
			beego.Error(err)
			return
		}

		username := u.Username
		//加载用户默认token
		tokenService := &service.TokenService{
			Session: db.GetMongoSession(),
		}
		defaultToken := tokenService.GetDefaultToken(username)
		if defaultToken != nil {
			this.Data["token"] = defaultToken.Value
		}
		//查询用户信息
		relation := user.Relation{}
		relationUser := relation.GetUserByName(username)
		this.Data["user"] = relationUser
		err = sess.Set(session.SESSION_USER, relationUser)
		if err != nil {
			beego.Error(err)
			return
		}
		// check if the code is right.
		err = sess.Set(session.SESSION_USERNAME, username)
		if err != nil {
			beego.Error(err)
			return
		}
		gitlab.Tokens.Add(username, access)

		beego.Info("Login ... ", username, "access token:", access.AccessToken)
		redirectUrl := sess.Get("redirect")
		if redirectUrl != nil {
			if r, ok := redirectUrl.(string); ok && r != "" {
				sess.Delete("redirect")
				http.Redirect(this.Ctx.ResponseWriter, this.Ctx.Request, r, 301)
				return
			}
		}
	} else {
		//全局模版变量
		this.Data["userName"] = sessUsername
		relation := user.Relation{}
		relationUser := relation.GetUserByName(sessUsername.(string))
		this.Data["user"] = relationUser
		//加载用户默认token
		tokenService := &service.TokenService{
			Session: db.GetMongoSession(),
		}
		defaultToken := tokenService.GetDefaultToken(sessUsername.(string))
		if defaultToken != nil {
			this.Data["token"] = defaultToken.Value
		}
		if name, ok := sessUsername.(string); ok {
			this.Username = name
		}
		// Already login
		//beego.Debug("Have code.", "sessCode",sessCode,"paramCode",paramCode)
	}
}
