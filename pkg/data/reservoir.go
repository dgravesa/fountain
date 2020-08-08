package data

import (
	"github.com/dgravesa/fountain/pkg/fountain"
)

// Reservoir is a data container for user water logs
type Reservoir interface {
	WriteWl(userID string, wl *fountain.WaterLog) error
	UserWls(userID string) ([]*fountain.WaterLog, error)
}
