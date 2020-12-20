package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/fountain"
	"github.com/dgravesa/minicli"
)

func main() {
	minicli.Func("user get", "print user logs", getUser)
	minicli.Func("user put", "create or update a user", putUser)
	minicli.Func("log put", "add a user log", putLog)

	if err := minicli.Exec(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getUser(args []string) error {
	var userID string

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&userID, "user", "", "ID of user")
	flags.Parse(args)

	if userID == "" {
		return fmt.Errorf("no user specified")
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
}

func putLog(args []string) error {
	var userID string
	var amt float64

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&userID, "user", "", "ID of user")
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
}

func putUser(args []string) error {
	var userID string

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&userID, "user", "", "ID of user")
	flags.Parse(args)

	if userID == "" {
		return fmt.Errorf("user not specified")
	}

	// try pulling user info
	userStore, err := data.DefaultUserStore()
	if err != nil {
		return err
	}
	user, err := userStore.User(userID)
	if err != nil {
		return err
	}

	// create new user with ID
	if user == nil {
		user = new(fountain.User)
		user.ID = userID
	}

	// update user via interactive prompts
	interactiveBuildUser(user)

	err = userStore.PutUser(user)
	if err != nil {
		return err
	}

	fmt.Println("user saved successfully")
	return nil
}

func interactiveBuildUser(user *fountain.User) {
	r := bufio.NewReader(os.Stdin)

	promptWithDefault := func(prompt, dflt string) string {
		// display prompt
		if dflt != "" {
			prompt += fmt.Sprintf(" (%s)", dflt)
		}
		fmt.Printf("%s: ", prompt)

		// get input from user
		inputln, _ := r.ReadString('\n')
		input := strings.TrimRight(inputln, "\n")

		if input == "" {
			return dflt
		}

		return input
	}

	user.FullName = promptWithDefault("Enter full name", user.FullName)

	for {
		email := promptWithDefault("Enter email", user.Email)

		if email != "" {
			// validate email string
			addrSplit := strings.Split(email, "@")
			if len(addrSplit) == 2 && strings.Contains(addrSplit[1], ".") {
				user.Email = email
				break
			}

			fmt.Println("Not a valid email address!")
		}
	}
}
