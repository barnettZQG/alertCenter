package session

import (
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
)

func init() {
	mangeConfig := &session.ManagerConfig{CookieLifeTime:3600, CookieName:"gosessionid", Gclifetime:3600, EnableSetCookie:true}
	GlobalSessions, _ = session.NewManager("memory", mangeConfig)
	go GlobalSessions.GC()
}

func GetSession(ctx *context.Context) (session.Store, error) {
	sess, err := GlobalSessions.SessionStart(ctx.ResponseWriter, ctx.Request)
	if err != nil {
		beego.Error("get session error:", err)
		return sess, err
	}
	return sess, nil
}