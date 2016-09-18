package gitlab

import (
	"github.com/astaxie/beego"
	"net/url"
	"strings"
	"alertCenter/models"
	"strconv"
)

func GetGitlabUrl() string {
	return strings.TrimSuffix(beego.AppConfig.String("Gitlab"), "/")
}

func GetAdminAccessToken() string {
	return beego.AppConfig.String("GitlabAccessToken")
}

// call back url
func GetCallBackUrl() string {
	return beego.AppConfig.String("GitlabCallBackUrl")
}

// call back urll encoding
func GetCallBackUrlEncode() string {
	u := GetCallBackUrl()
	return url.QueryEscape(u)
}

// clientId
func GetGitlabClientId() string {
	return beego.AppConfig.String("GitlabOAuthClientId")
}

// sercetId
func GetGitlabSercetId() string {
	return beego.AppConfig.String("GitlabOAuthSercet")
}

// redirect gitlab url
func GetGitlabOAuthUrl() string {
	return GetGitlabUrl() + "/oauth/authorize?response_type=code&client_id=" + GetGitlabClientId() + "&redirect_uri=" + GetCallBackUrlEncode()
}

func ConvertGitlabGroupToAlertModel(gitlab GitlabGroup) ( *models.Team) {
	team := &models.Team{}
	team.ID = strconv.Itoa(gitlab.Id)
	team.Name = gitlab.Name
	return team
}

func ConvertGitlabUserToAlertModel(gitlab GitlabUser) ( *models.User) {
	user := &models.User{}
	user.ID = strconv.Itoa(gitlab.Id)
	user.Name = gitlab.Username
	user.AvatarURL = gitlab.AvatarUrl
	user.Mail = gitlab.Email
	//user.Phone = gitlab.
	user.RealName = gitlab.Name
	return user
}

func ConvertGitlabUsers(gitlab []*GitlabUser) ([]*models.User) {
	users := make([]*models.User,len(gitlab))
	for i, u := range gitlab {
		users[i] = ConvertGitlabUserToAlertModel(*u)
	}
	return users
}