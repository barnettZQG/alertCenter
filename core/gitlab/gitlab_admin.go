package gitlab

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

func SearchUserByUsername(username string) (*GitlabUser, error) {
	url := GetGitlabUrl() + GetUser + "?username=" + username

	token := GetAdminAccessToken()
	resp, err := RequestGitlabWithToken(token, url, "GET", nil)
	if err != nil {
		beego.Error(err)

		return nil, err
	}

	user := &GitlabUser{}
	err = json.Unmarshal(resp, &user)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return user, nil
}