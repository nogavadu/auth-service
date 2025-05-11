package model

type Role struct {
	ID    uint64 `db:"id"`
	Name  string `db:"name"`
	Level int    `db:"level"`
}
