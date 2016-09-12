package gitlab

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
)

const (
	apiVersion = "/api/v3"
	getUser = "/users?per_page=100"
	getGroup = "/groups?per_page=100"
	getGroupUsers = "/groups/:id/members?per_page=100"
	currentUser = "/user"
)

func GitlabApi(method, url string, body []byte) ([]byte, error) {
	client := http.Client{}

	b := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		beego.Error(err.Error())
		return []byte{}, err
	}

	req.Header.Set("PRIVATE-TOKEN", GetAdminAccessToken())
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
	url := GetGitlabUrl() + apiVersion + currentUser
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
