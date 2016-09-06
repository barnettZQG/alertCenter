package models

type AlertReceive struct {
	Version           string            `json:"version"`
	Status            string            `json:"status"`
	Alerts            []Alert           `json:"alerts"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
}

type WeAlert struct {
	// ToUser  string            `json:"touser"`
	// ToParty string            `json:"toparty"`
	ToTag   string            `json:"totag"`
	MsgType string            `json:"msgtype"`
	AgentID int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

type WeTag struct {
	TagId   int    `json:"tagid"`
	TagName string `json:"tagname"`
}

type WeiTagResult struct {
	TagList []WeTag `json:"taglist"`
	ErrCode int     `json:"errcode"`
	ErrMsg  string  `json:"errmsg"`
}
