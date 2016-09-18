package gitlab

import "time"

type GitlabAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64 `json:"created_at"`
}

// save this in map
type Token struct {
	Username    string
	AccessToken string `json:"access_token"`
	Expire      time.Time
}

type GitlabUser struct {
	Id               int
	Name             string
	Username         string
	State            string
	CreatedAt        time.Time `json:"created_at"`
	AvatarUrl	 string `json:"avatar_url"`
	IsAdmin          bool `json:"is_admin"`
	Bio              interface{}
	Email            string
	ProjectsLimit    int `json:"projects_limit"`
	CurrentSignInAt  time.Time `json:"current_sign_in_at"`
	Identities       []interface{} `json:"identities"`
	CanCreateGroup   bool `json:"can_create_group"`
	CanCreateProject bool `json:"can_create_project"`
	Private_token    string `json:"private_token"`
}

type GitlabGroup struct {
	Id          int
	Name        string
	Path        string
	Description string
}

