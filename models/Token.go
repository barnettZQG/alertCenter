package models

import "time"

type Token struct {
	Value      string
	CreateTime time.Time
	Project    string
	UserName   string
}
