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

	mustNotErr := func(err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}

	// retrieve user logs
	reservoir, err := data.DefaultReservoir()
	mustNotErr(err)
	waterlogs, err := reservoir.UserWls(userID)
	mustNotErr(err)

	// print all user logs
	total := 0.0
	for _, wl := range waterlogs {
		fmt.Println(wl)
		total += wl.Amount
	}
	fmt.Println("Total amount:", total, "oz")
}
