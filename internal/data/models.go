package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users  UserModel
	Tokens TokenModel
	Books  BookModel
	Carts  CartModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db}, // initialize a new UserModel instance
		Tokens: TokenModel{DB: db},
		Books:  BookModel{DB: db},
		Carts:  CartModel{DB: db},
	}
}
