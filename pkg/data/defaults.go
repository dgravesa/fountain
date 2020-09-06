package data

import (
	"github.com/dgravesa/fountain/pkg/data/redis"
)

// DefaultReservoir returns a default reservoir instance
func DefaultReservoir() (Reservoir, error) {
	return redis.NewReservoir("")
}

// DefaultUserStore returns a default user store instance
func DefaultUserStore() (UserStore, error) {
	return redis.NewUserStore("")
}
