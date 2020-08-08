package fountain

import (
	"fmt"
	"time"
)

// WaterLog is the fundamental time/amount pair for logging water consumption
type WaterLog struct {
	time.Time `datastore:"time"`
	Amount    float64 `datastore:"amount"`
}

// WlNow returns a WaterLog corresponding to right now
func WlNow(amount float64) WaterLog {
	return WaterLog{Time: time.Now(), Amount: amount}
}

func (wl WaterLog) String() string {
	return fmt.Sprint(wl.Amount, " oz @ ", wl.Time)
}
