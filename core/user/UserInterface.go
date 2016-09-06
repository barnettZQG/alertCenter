package user

import "alertCenter/models"

type UserInterface interface {
	SearchTeams() ([]*models.Team, error)
	SearchUsers() ([]*models.User, error)
	GetUserByTeam(id string) ([]*models.User, error)
}
