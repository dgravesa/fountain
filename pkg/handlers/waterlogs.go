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
		writeResponse(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (waterlogs *WaterLogsHandler) get(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user")

	// validate user field
	if userID == "" {
		writeResponse(w, http.StatusBadRequest, "no user specified")
		return
	}

	// get user logs
	if userlogs, err := waterlogs.reservoir.UserWls(userID); err != nil {
		log.Println("internal error on user waterlogs get:", err)
		writeResponse(w, http.StatusInternalServerError, "internal server error")
	} else {
		writeResponse(w, http.StatusOK, userlogs)
	}
}

func (waterlogs *WaterLogsHandler) post(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user")
	amountStr := r.FormValue("amount")

	// validate user field
	if userID == "" {
		writeResponse(w, http.StatusBadRequest, "no user specified")
		return
	}

	// validate amount field
	if amount, err := strconv.ParseFloat(amountStr, 64); err != nil {
		writeResponse(w, http.StatusBadRequest, "unable to parse amount")
	} else if amount <= 0.0 {
		writeResponse(w, http.StatusBadRequest, "amount must be greater than 0.0")
	} else {
		// insert new log
		wl := fountain.WlNow(amount)
		if err := waterlogs.reservoir.WriteWl(userID, &wl); err != nil {
			log.Println("internal error on user waterlog post:", err)
			writeResponse(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// success response
		writeResponse(w, http.StatusOK, wl)
	}
}
