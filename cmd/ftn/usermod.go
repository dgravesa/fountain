package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dgravesa/fountain/pkg/data"
	"github.com/dgravesa/fountain/pkg/fountain"
)

func modifyUser(args []string) error {
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
