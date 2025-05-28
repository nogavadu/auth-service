package model

type UserInfo struct {
	Id     int    `json:"id"`
	Email  string `json:"email"`
	RoleId int    `json:"role"`
}
