package model

type UserInfo struct {
	Email  string `json:"email"`
	RoleId uint64 `json:"role"`
}
