package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/fountain"
	"github.com/dgravesa/minicli"
)

var userID string

func main() {
	minicli.Flags("", "", func(flags *flag.FlagSet) {
		flags.StringVar(&userID, "user", "", "ID of user")
	})

	minicli.Cmd("serve", "start the fountain service", &serveCmd{})

	minicli.Func("get", "get user logs", func(args []string) error {
		hasErrors := false

		if userID == "" {
			log.Println("no user specified")
			hasErrors = true
		}

		if hasErrors {
			log.Fatalln("error occurred")
		}

		// retrieve user logs
		reservoir, err := data.DefaultReservoir()
		if err != nil {
			return err
		}
		waterlogs, err := reservoir.UserWls(userID)
		if err != nil {
			return err
		}

		// print all user logs
		total := 0.0
		for _, wl := range waterlogs {
			fmt.Println(wl)
			total += wl.Amount
		}
		fmt.Println("Total amount:", total, "oz")

		return nil
	})

	minicli.Func("put", "put a new user log", func(args []string) error {
		var amt float64
		flags := flag.NewFlagSet("put", flag.ExitOnError)
		flags.Float64Var(&amt, "amount", 0.0, "amount of water")
		flags.Parse(args)

		if userID == "" {
			return fmt.Errorf("no user specified")
		}

		if amt <= 0.0 {
			return fmt.Errorf("amount must be greater than 0.0")
		}

		// insert new user log
		wl := fountain.WlNow(amt)
		reservoir, err := data.DefaultReservoir()
		if err != nil {
			return err
		}
		err = reservoir.WriteWl(userID, &wl)
		if err != nil {
			return err
		}
		fmt.Println(wl)
		return nil
	})

	minicli.Func("usermod", "modify a fountain user", modifyUser)

	if err := minicli.Exec(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
