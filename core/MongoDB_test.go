package core

import (
	"testing"

	"github.com/barnettZQG/alertCenter/util"
)

type Person struct {
	NAME  string
	PHONE string
}

func Test_Insert(t *testing.T) {
	session := GetMongoSession()
	if ok := session.Insert("Person", &Person{PHONE: "18811577546",
		NAME: "barnett"}); ok {
		util.Info("insert success")
	}

	session.Close()
}
