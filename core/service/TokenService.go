package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type TokenService struct {
	Session *db.MongoSession
}

var cacheToken map[string][]*models.Token

func init() {
	if cacheToken == nil {
		cacheToken = make(map[string][]*models.Token, 0)
		Session := db.GetMongoSession()
		col := Session.GetCollection("Token")
		if col == nil {
			beego.Error("get collection token error ")
		}
		var result []*models.Token
		err := col.Find(nil).All(&result)
		if err != nil {
			beego.Error("find all token error." + err.Error())
		}
		for _, item := range result {
			items := cacheToken[item.UserName]
			if items == nil {
				items = make([]*models.Token, 0)
			}
			items = append(items, item)
			cacheToken[item.UserName] = items
		}
	}
}

//GetDefaultToken 获取默认token
func (e *TokenService) GetDefaultToken(user string) *models.Token {

	tokens := cacheToken[user]
	for _, token := range tokens {
		if token.Project == "default" {
			return token
		}
	}
	return nil
}

//CreateToken 创建token
func (e *TokenService) CreateToken(project string, userName string) *models.Token {
	beego.Debug("create token project:" + project + " ,user:" + userName)
	token := &models.Token{
		Value:      uuid.NewV4().String(),
		CreateTime: time.Now(),
		Project:    project,
		UserName:   userName,
	}
	items := cacheToken[userName]
	if items == nil {
		items = make([]*models.Token, 0)
	}
	items = append(items, token)
	cacheToken[userName] = items
	e.Session.Insert("Token", token)
	return token
}

//CheckToken 验证token
func (e *TokenService) CheckToken(token string, user string) bool {
	start := time.Now()
	defer fmt.Println("checkToken time:", time.Now().Sub(start))
	for _, v := range cacheToken[user] {
		if v.Value == token {
			return true
		}
	}
	return false
}

//DeleteToken 删除token
func (e *TokenService) DeleteToken(project string, user string) bool {
	col := e.Session.GetCollection("Token")
	if col == nil {
		beego.Error("get collection token error ")
		return false
	}
	err := col.Remove(bson.M{"project": project, "username": user})
	if err != nil {
		beego.Error("delete token error ,", err.Error())
		return false
	}
	//删除slice成员
	var items []*models.Token
	for _, v := range cacheToken[user] {
		if v.Project != project {
			items = append(items, v)
		}
	}
	cacheToken[user] = items
	return true
}

//GetToken 获取token
func (e *TokenService) GetToken(project string, user string) *models.Token {
	tokens := cacheToken[user]
	for _, token := range tokens {
		if token.Project == project {
			return token
		}
	}
	return nil
}

//GetAllToken 获取用户所有token
func (e *TokenService) GetAllToken(userName string) (result []*models.Token) {
	for _, v := range cacheToken[userName] {
		if v.Project != "default" {
			result = append(result, v)
		}
	}
	return
}
