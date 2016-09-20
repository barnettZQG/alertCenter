package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

const (
	ApiVersion       = "/api/v3"
	GetUser          = "/users"
	GetGroups        = "/groups"
	GetGroupUsers    = "/groups/:id/members"
	currentUser      = "/user"
	SearchUserByName = "/users?username=:username"
)

func GitlabApi(method, url string, body []byte) ([]byte, error) {
	client := http.Client{}

	b := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		beego.Error(err.Error())
		return []byte{}, err
	}

	//if way == ""
	req.Header.Set("PRIVATE-TOKEN", GetAdminAccessToken())
	//fmt.Println("debug:", "private token:", GetAdminAccessToken(), "url:", url)
	//req.Header.Set("Authorization", "Bearer " + accessToken)
	resp, err := client.Do(req)
	if err != nil {
		beego.Error(err.Error())
		return []byte{}, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err.Error())
		return []byte{}, err
	}
	defer resp.Body.Close()
	return respBody, nil
}

func GetCurrentUserWithToken(token string) (*GitlabUser, error) {
	url := GetGitlabUrl() + ApiVersion + currentUser
	body, err := RequestGitlabWithToken(token, url, "GET", nil)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	user := &GitlabUser{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*GitlabUser, error) {
	userUrl := strings.Replace(SearchUserByName, ":username", username, -1)

	url := GetGitlabUrl() + ApiVersion + userUrl
	//fmt.Println("url:", url)
	body, err := GitlabApi("GET", url, nil)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	user := []*GitlabUser{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		beego.Debug("Unmarshal:", string(body))
		beego.Error(err)
		return nil, err
	}
	if len(user) == 1 {
		return user[0], nil
	} else {
		return nil, fmt.Errorf("Search more then one user.")
	}
}

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
