package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/dgravesa/fountain/pkg/fountain"
)

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

func main() {
	var userID string

	flag.StringVar(&userID, "user", "", "ID of user")
	flag.Parse()

	if userID == "" {
		log.Fatalln("user not specified")
	}

	// initialize datastore client
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "water-you-logging-for")
	if err != nil {
		log.Fatalln(err)
	}

	// try pulling user info from datastore
	user := fountain.User{ID: userID}
	k := datastore.NameKey("Users", userID, nil)
	client.Get(ctx, k, &user)

	// update user via interactive prompts
	interactiveBuildUser(&user)

	// insert updated user into datastore
	if _, err = client.Put(ctx, k, &user); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("user saved successfully")
}
