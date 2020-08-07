package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

type waterlog struct {
	Amount float64
	time.Time
}

func main() {
	var user string
	var amt float64

	flag.StringVar(&user, "user", "", "name of user")
	flag.Float64Var(&amt, "amount", 0.0, "amount of water")
	flag.Parse()

	hasErrors := false

	if user == "" {
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

	// insert new item
	k := datastore.NameKey("Its Me", user, nil)
	wl := waterlog{Amount: amt, Time: time.Now()}
	if _, err := client.Put(ctx, k, &wl); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("successfully made it this far")
}
