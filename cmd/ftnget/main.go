package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/api/iterator"

	"github.com/dgravesa/fountain/pkg/fountain"

	"cloud.google.com/go/datastore"
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

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, "water-you-logging-for")
	if err != nil {
		log.Fatalln(err)
	}

	userKey := datastore.NameKey("Users", userID, nil)
	q := datastore.NewQuery("WaterLogs").Ancestor(userKey)
	qResult := client.Run(ctx, q)

	for {
		var wl fountain.WaterLog

		// retrieve next log
		if _, err = qResult.Next(&wl); err != nil {
			if err == iterator.Done {
				break
			} else {
				log.Fatalln(err)
			}
		}

		fmt.Println(wl.Amount, "oz @", wl.Time)
	}
}
