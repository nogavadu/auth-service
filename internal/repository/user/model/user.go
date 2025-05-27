package model

type User struct {
	ID       uint64 `db:"id"`
	Email    string `db:"email"`
	PassHash string `db:"password_hash"`
	RoleId   uint64 `db:"role"`
}
