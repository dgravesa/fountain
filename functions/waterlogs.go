package functions

import (
	"net/http"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/handlers"
)

// ServeWaterLogs is a HandleFunc for the waterlogs resource
func ServeWaterLogs(w http.ResponseWriter, r *http.Request) {
	reservoir := data.DefaultReservoir()
	handler := handlers.NewWaterLogsHandler(reservoir)
	handler.ServeHTTP(w, r)
}
