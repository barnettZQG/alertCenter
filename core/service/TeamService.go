package service

import (
	"alertCenter/core/db"
	"alertCenter/models"
)

type TeamService struct {
	Session *db.MongoSession
}

func GetTeamService(session *db.MongoSession) *TeamService {
	return &TeamService{
		Session: session,
	}
}

func (e *TeamService) FindAll() (teams []*models.Team) {
	coll := e.Session.GetCollection("team")
	if coll == nil {
		return nil
	}
	coll.Find(nil).Select(nil).All(&teams)
	return
}
