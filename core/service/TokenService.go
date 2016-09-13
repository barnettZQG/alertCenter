package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type TokenService struct {
	Session *db.MongoSession
}

var cacheToken map[string]*models.Token

func init() {
	tokenService := &TokenService{
		Session: db.GetMongoSession(),
	}
	tokenService.GetToken("")
}

//CreateToken 创建token
func (e *TokenService) CreateToken(project string, userName string) *models.Token {
	token := &models.Token{
		Value:      uuid.NewV4().String(),
		CreateTime: time.Now(),
		Project:    project,
		UserName:   userName,
	}
	cacheToken[project] = token
	e.Session.Insert("Token", token)
	return token
}

//CheckToken 验证token
func (e *TokenService) CheckToken(token string, user string) bool {
	for _, v := range cacheToken {
		if v.Value == token && v.UserName == user {
			return true
		}
	}
	return false
}

//DeleteToken 删除token
func (e *TokenService) DeleteToken(object string) bool {
	col := e.Session.GetCollection("Token")
	if col == nil {
		beego.Error("get collection token error ")
		return false
	}
	err := col.Remove(bson.M{"object": object})
	if err != nil {
		beego.Error("delete token error ,", err.Error())
		return false
	}
	delete(cacheToken, object)
	return true
}

//GetToken 获取token
func (e *TokenService) GetToken(object string) *models.Token {
	if cacheToken == nil {
		cacheToken = make(map[string]*models.Token, 0)
		col := e.Session.GetCollection("Token")
		if col == nil {
			beego.Error("get collection token error ")
			return nil
		}
		var result []*models.Token
		err := col.Find(nil).All(&result)
		if err != nil {
			beego.Error("find all token error." + err.Error())
		}
		for _, item := range result {
			cacheToken[item.Project] = item
		}
		return cacheToken[object]
	}
	return cacheToken[object]
}

//GetAllToken 获取用户所有token
func (e *TokenService) GetAllToken(userName string) (result []*models.Token) {
	for _, v := range cacheToken {
		if v.UserName == userName {
			result = append(result, v)
		}
	}
	return
}
