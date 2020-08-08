package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
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

	ctx := context.Background()

	// initialize client
	client, err := datastore.NewClient(ctx, "water-you-logging-for")
	if err != nil {
		log.Fatalln(err)
	}

	// create log
	wl := fountain.WlNow(amt)

	// insert new item
	userKey := datastore.NameKey("Users", userID, nil)
	wlKey := datastore.IDKey("WaterLogs", wl.Unix(), userKey)
	if _, err := client.Put(ctx, wlKey, &wl); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("success")
}
