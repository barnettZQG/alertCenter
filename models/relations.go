package models

type Receiver struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Users     []*User  `json:"mail"`
	WeGroupID []string `json:"weiGroupId"`
}
type Team struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	WeTeamID     string `json:"weId"`
	ParentTeamID string `json:"parentId"`
}

type APP struct {
	ID           string `json:"app_id"`
	Name         string `json:"app_name"`
	Users        []*User
	Teams        []*Team
	Mails        []string `json:"emails"`
	IDC          string   `json:"idc"`
	Domain       string   `json:"domain"`
	BusinessLine string   `json:"business_line"`
}

type User struct {
	ID     string `json:"id"`
	Name   string `json:"Name"`
	TeamID string `json:"teamId"`
	Phone  string `json:"phone"`
	Mail   string `json:"mail"`
	WeID   string `json:"weId"`
}
