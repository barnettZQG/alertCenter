package session

import (
	"github.com/astaxie/beego/session"
)

var GlobalSessions *session.Manager

const (
	SESSION_USER = "user"
	SESSION_USERNAME = "username"
)
