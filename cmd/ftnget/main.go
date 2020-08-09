package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dgravesa/fountain/pkg/data"
)

func main() {
	var userID string

	flag.StringVar(&userID, "user", "", "ID of user")
	flag.Parse()

	hasErrors := false

	if userID == "" {
		log.Println("no user specified")
		hasErrors = true
	}

	if hasErrors {
		log.Fatalln("error occurred")
	}

	// retrieve user logs
	reservoir := data.DefaultReservoir()
	waterlogs, err := reservoir.UserWls(userID)

	if err != nil {
		log.Fatalln(err)
	} else {
		total := 0.0

		// print all user logs
		for _, wl := range waterlogs {
			fmt.Println(wl)
			total += wl.Amount
		}

		fmt.Println("Total amount:", total, "oz")
	}
}
