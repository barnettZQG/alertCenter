package gitlab

import (
	"github.com/astaxie/beego"
	"net/url"
	"net/http"
	"io/ioutil"
	"bytes"
	"strings"
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

func RequestGitlab(username, url, method string, body []byte) ([]byte, error) {
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
	req.Header.Add("Authorization", "Bearer " + token)

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


