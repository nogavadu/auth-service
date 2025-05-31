package model

type User struct {
	Id int `db:"id"`
	UserInfo
}

type UserInfo struct {
	Name     *string `db:"name"`
	Email    string  `db:"email"`
	PassHash string  `db:"password_hash"`
	Avatar   *string `db:"avatar"`
	RoleId   int     `db:"role"`
}

type UserUpdateInput struct {
	Name     *string `db:"name"`
	Email    *string `db:"email"`
	Password *string `db:"password"`
	Avatar   *string `db:"avatar"`
	RoleId   *int    `db:"role"`
}
