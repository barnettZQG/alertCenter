package core

import "alertCenter/models"

type TeamService struct {
	Session *MongoSession
}

func GetTeamService(session *MongoSession) *TeamService {
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
