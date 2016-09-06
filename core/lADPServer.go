package core

import (
	"fmt"

	"alertCenter/models"

	"github.com/astaxie/beego"

	ldap "gopkg.in/ldap.v2"
)

var ldapServer string
var ldapDN string
var ldapPass string
var ldapPort int
var err error

func init() {
	ldapServer = beego.AppConfig.String("LADPServer")
	ldapPort, err = beego.AppConfig.Int("LDAPPort")
	ldapDN = beego.AppConfig.String("LDAPDN")
	ldapPass = beego.AppConfig.String("LDAPPass")
}

type LDAPServer struct {
}

func (e *LDAPServer) SearchTeams() (teams []*models.Team, err error) {
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	defer l.Close()
	err = l.Bind(ldapDN, ldapPass)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		"dc=yunpro,dc=cn", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=posixGroup))", // The filter to apply
		nil, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	for _, entry := range sr.Entries {
		team := &models.Team{
			ID:   entry.GetAttributeValue("gidNumber"),
			Name: entry.GetAttributeValue("cn"),
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (e *LDAPServer) SearchUsers() (users []*models.User, err error) {
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	defer l.Close()
	err = l.Bind(ldapDN, ldapPass)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		"dc=yunpro,dc=cn", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=posixAccount))", // The filter to apply
		nil, // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	for _, entry := range sr.Entries {
		user := &models.User{
			ID:     entry.GetAttributeValue("uidNumber"),
			Name:   entry.GetAttributeValue("cn"),
			TeamID: entry.GetAttributeValue("gidNumber"),
			Phone:  entry.GetAttributeValue("mobile"),
			Mail:   entry.GetAttributeValue("Email"),
		}
		users = append(users, user)
	}
	return users, nil
}
