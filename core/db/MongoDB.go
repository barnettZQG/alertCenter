package db

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

var Session *mgo.Session

func init() {
	session := createSession()
	if session == nil {
		panic("mongodb init error!")
	}
	Session = session
}
func createSession() *mgo.Session {
	URL := beego.AppConfig.String("mongoURI")
	fmt.Println("Url:", URL)
	//URL := "10.12.1.129:27017"
	session, err := mgo.Dial(URL) //连接数据库
	if err != nil {
		beego.Error("Get mongo session error." + err.Error())
		return nil
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	return session
}
func GetSession() (*mgo.Session, error) {
	if Session == nil {
		return createSession(), nil
	}
	return Session.Clone(), nil
}

func CloseSession(session *mgo.Session) {
	if session != nil {
		session.Close()
	}
}

func GetDB(session *mgo.Session) (*mgo.Database, error) {
	dbName := beego.AppConfig.String("mongoDB")
	dbUser := beego.AppConfig.String("mongoUser")
	dbPass := beego.AppConfig.String("mongoPass")
	// dbName := "admin"
	// dbUser := "root"
	// dbPass := "root"
	if session == nil {
		var err error
		session, err = GetSession()
		if err != nil {
			return nil, err
		}
	}
	db := session.DB(dbName)
	err := db.Login(dbUser, dbPass)
	if err != nil {
		beego.Error("Login mongodb userName or password error." + err.Error())
		return nil, err
	}
	return db, nil
}

func GetCollection(collection string, db *mgo.Database) (*mgo.Collection, error) {
	if len(collection) == 0 {
		return nil, errors.New("Don't use empty collection name")
	}
	if db == nil {
		var err error
		db, err = GetDB(nil)
		if err != nil {
			return nil, err
		}
	}
	co := db.C(collection)
	return co, nil
}

type MongoSession struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func GetMongoSession() *MongoSession {
	session, err := GetSession()
	if err != nil {
		return nil
	}
	db, err := GetDB(session)
	if err != nil {
		return nil
	}
	return &MongoSession{
		Session: session,
		DB:      db,
	}
}

func (e *MongoSession) Insert(collection string, data ...interface{}) bool {
	coll, err := GetCollection(collection, e.DB)
	if err != nil {
		return false
	}
	err = coll.Insert(data...)
	if err != nil {
		beego.Error("insert data in " + collection + " error," + err.Error())
		return false
	}
	return true
}

func (e *MongoSession) GetCollection(collection string) *mgo.Collection {
	coll, err := GetCollection(collection, e.DB)
	if err != nil {
		return nil
	}
	return coll
}
func (e *MongoSession) RemoveAll(collection string) bool {
	coll, err := GetCollection(collection, e.DB)
	if err != nil {
		return false
	}
	_, err = coll.RemoveAll(nil)
	if err == nil {
		return true
	}
	return false
}

func (e *MongoSession) Close() {
	if e.Session != nil {
		CloseSession(e.Session)
	}
}
