package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/fountain"
)

// WaterLogsHandler is a Handler for requests to the waterlogs resource
type WaterLogsHandler struct {
	reservoir data.Reservoir
}

// NewWaterLogsHandler instantiates a new handler for the waterlogs resource
func NewWaterLogsHandler(r data.Reservoir) *WaterLogsHandler {
	h := new(WaterLogsHandler)
	h.reservoir = r
	return h
}

func (waterlogs *WaterLogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		waterlogs.get(w, r)
	case http.MethodPost:
		waterlogs.post(w, r)
	default:
		errors := []string{"method not allowed"}
		writeResponse(w, http.StatusMethodNotAllowed, errors)
	}
}

func (waterlogs *WaterLogsHandler) get(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	errors := []string{"not yet implemented"}
	writeResponse(w, http.StatusNotImplemented, errors)
}

func (waterlogs *WaterLogsHandler) post(w http.ResponseWriter, r *http.Request) {
	var errors []string
	statusCode := http.StatusOK

	userID := r.FormValue("user")
	amountStr := r.FormValue("amount")

	// validate user field
	if userID == "" {
		errors = append(errors, "no user specified")
		statusCode = http.StatusBadRequest
	}

	if len(errors) == 0 {
		amount, err := strconv.ParseFloat(amountStr, 64)

		// validate amount field
		if err != nil {
			errors = append(errors, "unable to parse amount")
			statusCode = http.StatusBadRequest
		} else if amount <= 0.0 {
			errors = append(errors, "amount must be greater than 0.0")
			statusCode = http.StatusBadRequest
		} else {
			// insert new log
			wl := fountain.WlNow(amount)
			if err := waterlogs.reservoir.WriteWl(userID, &wl); err != nil {
				log.Println("internal error on user waterlog post:", err)
				statusCode = http.StatusInternalServerError
			}
		}
	}

	if statusCode != http.StatusOK {
		log.Println("errors on user waterlog post:", errors)
	}

	writeResponse(w, statusCode, errors)
}
