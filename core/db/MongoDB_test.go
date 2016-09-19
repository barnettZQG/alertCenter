package db

import "testing"

type Person struct {
	NAME  string
	PHONE string
}

func Test_Insert(t *testing.T) {
	session := GetMongoSession()
	if session != nil {
		defer session.Close()
	}
	if ok := session.Insert("Person", &Person{PHONE: "18811577546",
		NAME: "barnett"}); ok {
	}

	session.Close()
}
