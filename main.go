package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          int             `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
}

type UserActivity struct {
	Balance      float64       `json:"balance"`
	Password     string        `json:"password"`
	Transactions []Transaction `json:"transactions"`
}

func login(user string, password string) (UserActivity, error) {
	filePath := "./data/" + user + ".json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		return UserActivity{}, fmt.Errorf("file not found: %s", filePath)
	}

	var userActivity UserActivity
	err = json.Unmarshal(data, &userActivity)
	if err != nil {
		return UserActivity{}, fmt.Errorf("error parsing file: %v", err)
	}

	if password != userActivity.Password {
		return UserActivity{}, fmt.Errorf("incorrect password")
	}

	return userActivity, nil
}

func main() {
	fmt.Println("Welcome to my expense tracker")
	fmt.Println("=============================")
	fmt.Println("Add User (1)")
	fmt.Println("Login User (2)")

	for {
		var choice int8
		fmt.Print("\nChoose an option: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Println("Add User (not implemented yet)")
		case 2:
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Your username: ")
			user, _ := reader.ReadString('\n')
			user = strings.TrimSpace(user)

			fmt.Print("Your password: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			userActivity, err := login(user, password)
			if err != nil {
				fmt.Println("Login failed:", err)
				continue
			}

			fmt.Println("Login successful!")
			fmt.Printf("Welcome, %s! Your current balance is: $%.2f\n", user, userActivity.Balance)

		default:
			fmt.Println("No valid Option")
		}
	}
}
