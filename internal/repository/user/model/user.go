package model

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	PassHash string `db:"password_hash"`
	RoleId   int    `db:"role"`
}
