package data

import (
	"github.com/dgravesa/fountain/pkg/fountain"
)

// UserStore is a data container for Fountain users
type UserStore interface {
	User(userID string) (*fountain.User, error)
	PutUser(user *fountain.User) error
}
