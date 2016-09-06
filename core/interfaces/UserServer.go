package interfaces

import "alertCenter/models"

type UserManager interface {
	SearchTeams() (teams []*models.Team, err error)
	SearchUsers() (users []*models.User, err error)
}
