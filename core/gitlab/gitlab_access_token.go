package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego"
)

var Tokens GitlabTokens

//统一用username作为key,因为gitlab username是唯一的,name是可以重复的
type GitlabTokens map[string]Token

func init() {
	Tokens = make(map[string]Token)
}

func (t GitlabTokens) Add(username string, at *GitlabAccessToken) {
	token := Token{AccessToken: at.AccessToken, Username: username}
	token.Expire = time.Unix(at.CreatedAt, 0).Add(time.Hour * 2)
	t[username] = token
}

func (t GitlabTokens) Delete(username string) {
	delete(t, username)
}

func (t GitlabTokens) Get(username string) (Token, error) {
	token, ok := t[username]
	if !ok {
		return token, fmt.Errorf("Token not found")
	}

	if token.Expire.Before(time.Now()) {
		return token, fmt.Errorf("Token is expired.")
	}

	return token, nil
}

func (t GitlabTokens) Update(username string, at *GitlabAccessToken) {
	t.Delete(username)
	t.Add(username, at)
}

// get access token
func GetGitlabAccessToken(code string) (*GitlabAccessToken, error) {
	u, err := url.Parse("/oauth/token")
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	query := u.Query()
	query.Add("client_id", GetGitlabClientId())
	query.Add("client_secret", GetGitlabSercetId())
	query.Add("code", code)
	query.Add("grant_type", "authorization_code")
	query.Add("redirect_uri", GetCallBackUrl())

	q := query.Encode()
	uri := u.Path + "?" + q

	targetUrl := beego.AppConfig.String("Gitlab") + uri
	beego.Debug("GetGitlabAccessToken url:", targetUrl, "gitlab:", beego.AppConfig.String("Gitlab"))

	resp, err := http.Post(targetUrl, "application/json", nil)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	accToken := &GitlabAccessToken{}
	err = json.Unmarshal(body, &accToken)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	//fmt.Printf("body is %#v\n", string(body))

	return accToken, nil
}
