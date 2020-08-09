package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dgravesa/fountain/pkg/data"

	"github.com/dgravesa/fountain/pkg/fountain"
)

func main() {
	var userID string
	var amt float64

	flag.StringVar(&userID, "user", "", "ID of user")
	flag.Float64Var(&amt, "amount", 0.0, "amount of water")
	flag.Parse()

	hasErrors := false

	if userID == "" {
		log.Println("no user specified")
		hasErrors = true
	}

	if amt <= 0.0 {
		log.Println("amount must be greater than 0.0")
		hasErrors = true
	}

	if hasErrors {
		log.Fatalln("error occurred")
	}

	// insert new user log
	wl := fountain.WlNow(amt)
	reservoir := data.DefaultReservoir()
	err := reservoir.WriteWl(userID, &wl)

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(wl)
	}
}
