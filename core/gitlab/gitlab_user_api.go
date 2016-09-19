package gitlab

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

func RequestGitlab(username, method, url string, body []byte) ([]byte, error) {
	token, err := Tokens.Get(username)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return RequestGitlabWithToken(token.AccessToken, url, method, body)
}

func RequestGitlabWithToken(token, url, method string, body []byte) ([]byte, error) {
	client := http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		beego.Error(err)
		return nil, nil
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	defer resp.Body.Close()

	return b, nil
}

func GetGroupsByUsername(username string) ([]GitlabGroup, error) {
	url := GetGitlabUrl() + ApiVersion + GetGroups
	//fmt.Println("url:", url)
	body, err := RequestGitlab(username, "GET", url, nil)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	groups := []GitlabGroup{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		beego.Debug("Unmarshal:", string(body))
		beego.Error(err)
		return nil, err
	}

	return groups, nil
}

func GetUsersByTeam(username, groupid string) ([]*GitlabUser, error) {
	gurl := strings.Replace(GetGroupUsers, ":id", groupid, -1)
	url := GetGitlabUrl() + ApiVersion + gurl
	//fmt.Println("url:", url)
	body, err := RequestGitlab(username, "GET", url, nil)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	users := []*GitlabUser{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		beego.Debug("Unmarshal:", string(body))
		beego.Error(err)
		return nil, err
	}

	return users, nil
}
