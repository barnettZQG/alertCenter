package user

import "alertCenter/models"

type UserInterface interface {
	//查询所有团队
	SearchTeams() ([]*models.Team, error)
	//查询所有用户
	SearchUsers() ([]*models.User, error)
	//根据团队查询用户
	GetUserByTeam(id string) ([]*models.User, error)
	//根据名称获取用户
	//GetUserByUserName() (*models.User, error)
}

//UserAuthentication 用户管理之用户权限认真
type UserAuthentication interface {

	//API权限验证方法
	APIFilter()
	//页面请求权限验证方法
	HTTPFilter()
}
