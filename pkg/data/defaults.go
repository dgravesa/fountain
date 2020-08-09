package data

import (
	"github.com/dgravesa/fountain/pkg/data/gcp"
)

// DefaultReservoir returns a default reservoir instance
func DefaultReservoir() Reservoir {
	return gcp.DatastoreClient{}
}

// DefaultUserStore returns a default user store instance
func DefaultUserStore() UserStore {
	return gcp.DatastoreClient{}
}
