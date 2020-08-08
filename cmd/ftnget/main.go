package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/dgravesa/fountain/pkg/waterlog"
)

func main() {
	var user string

	flag.StringVar(&user, "user", "", "name of user")
	flag.Parse()

	hasErrors := false

	if user == "" {
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

	k := datastore.NameKey("Its Me", user, nil)
	var wl waterlog.WaterLog
	if err = client.Get(ctx, k, &wl); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(wl.Amount, "oz @", wl.Time)
}
