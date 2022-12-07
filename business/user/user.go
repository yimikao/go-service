package user

import "github.com/go-pg/pg"

type Storer interface {
	GetById(id int64) (*User, error)
	All() ([]*User, error)
	Create(cm *User) (*User, error)
}

type userLayer struct {
	dbConn *pg.DB
}
