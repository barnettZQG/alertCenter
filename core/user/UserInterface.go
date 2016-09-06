package user

import "alertCenter/models"

type UserInterface interface {
	SearchTeams() ([]*models.Team, error)
	SearchUsers() ([]*models.User, error)
}
