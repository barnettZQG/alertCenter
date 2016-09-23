package models

import (
	"time"

	"github.com/prometheus/common/model"
	"gopkg.in/mgo.v2/bson"
)

type Alert struct {
	Labels        Label     `json:"labels"`
	Annotations   Label     `json:"annotations"`
	StartsAt      time.Time `json:"startsAt,omitempty"`
	EndsAt        time.Time `json:"endsAt,omitempty"`
	GeneratorURL  string    `json:"generatorURL"`
	Mark          string    `json:"mark" bson:"mark"`
	Receiver      *Receiver `json:"receiver"`
	AlertCount    int
	IsHandle      int
	HandleDate    time.Time `json:"handleDate,omitempty"`
	HandleMessage string
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
	Level         int       `json:"level,omitempty"`
}

type Label struct {
	model.LabelSet
}

//Contains a是否包含source
func (a Label) Contains(source Label) bool {
	for k, v := range source.LabelSet {
		if va, ok := a.LabelSet[k]; ok {
			if v != va {
				return false
			}
		} else {
			return false
		}
	}
	return true
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
	res.Annotations = o.Annotations
	return &res
}

//Reset 重置alert状态
func (a *Alert) Reset(o *Alert) *Alert {
	res := *a
	res.StartsAt = o.StartsAt
	res.EndsAt = o.EndsAt
	res.Annotations = o.Annotations
	if res.EndsAt.IsZero() {
		res.AlertCount = 1
		res.IsHandle = 0
		res.HandleDate = time.Now()
		res.HandleMessage = "报警再次产生"
	} else {
		res.IsHandle = 2
		res.HandleDate = time.Now()
		res.HandleMessage = "报警已自动恢复"
	}
	res.UpdatedAt = time.Now()
	return &res
}

type AlertHistory struct {
	ID       bson.ObjectId `bson:"_id"`
	Mark     string        `json:"mark"`
	AddTime  time.Time     `json:"addTime"`
	StartsAt time.Time     `json:"startsat"`
	EndsAt   time.Time     `json:"endsat"`
	Duration time.Duration `json:"duration"`
	Message  string        `json:"message"`
	Value    string        `json:"value"`
}
