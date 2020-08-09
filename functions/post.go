package functions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/fountain"
)

// FountainPost handles the a WaterLog POST request
func FountainPost(w http.ResponseWriter, r *http.Request) {
	// TODO: this will need to handle both get and post methods, or the endpoint renamed
	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, []string{"method not allowed"})
		return
	}

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
			if err := data.DefaultReservoir().WriteWl(userID, &wl); err != nil {
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
