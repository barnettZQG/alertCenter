package models

import (
	"time"

	"github.com/prometheus/common/model"
	"gopkg.in/mgo.v2/bson"
)

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

type Alert struct {
	_ID           bson.ObjectId
	Labels        model.LabelSet `json:"labels"`
	Annotations   model.LabelSet `json:"annotations"`
	StartsAt      time.Time      `json:"startsAt,omitempty"`
	EndsAt        time.Time      `json:"endsAt,omitempty"`
	GeneratorURL  string         `json:"generatorURL"`
	Mark          string         `json:"mark" bson:"mark"`
	Receiver      *Receiver      `json:"receiver"`
	AlertCount    int
	IsHandle      int
	HandleDate    time.Time `json:"handleDate,omitempty"`
	HandleMessage string
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

func (a *Alert) Fingerprint() model.Fingerprint {
	return a.Labels.Fingerprint()
}

// Merge merges the timespan of two alerts based and overwrites annotations
// based on the authoritative timestamp.  A new alert is returned, the labels
// are assumed to be equal.
func (a *Alert) Merge(o *Alert) *Alert {
	// Let o always be the younger alert.
	if o.UpdatedAt.Before(a.UpdatedAt) {
		return o.Merge(a)
	}

	res := *a

	// Always pick the earliest starting time.
	if a.StartsAt.After(o.StartsAt) {
		res.StartsAt = o.StartsAt
	}

	// A non-timeout resolved timestamp always rules.
	// The latest explicit resolved timestamp wins.
	if a.EndsAt.Before(o.EndsAt) {
		res.EndsAt = o.EndsAt
	}
	return &res
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

