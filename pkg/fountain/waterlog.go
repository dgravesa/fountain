package fountain

import "time"

// WaterLog is the fundamental time/amount pair for logging water consumption
type WaterLog struct {
	time.Time
	Amount float64
}

// WlNow returns a WaterLog corresponding to right now
func WlNow(amount float64) WaterLog {
	return WaterLog{Time: time.Now(), Amount: amount}
}
