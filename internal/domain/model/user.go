package model

type User struct {
	Id int `json:"id"`
	UserInfo
}

type UserInfo struct {
	Name   *string `json:"name"`
	Email  string  `json:"email"`
	Avatar *string `json:"avatar"`
	Role   string  `json:"role"`
}

type UserUpdateInput struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"roleId,omitempty"`
}
