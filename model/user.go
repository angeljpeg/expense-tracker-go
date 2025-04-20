package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type UserActivity struct {
	Balance      float64       `json:"balance"`
	Password     string        `json:"password"`
	Transactions []Transaction `json:"transactions"`
}

func login(user string, password string) (UserActivity, error) {
	filePath := "./data/" + user + ".json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return UserActivity{}, fmt.Errorf("user not found")
	}
	var ua UserActivity
	if err := json.Unmarshal(data, &ua); err != nil {
		return UserActivity{}, fmt.Errorf("error parsing file")
	}
	if password != ua.Password {
		return UserActivity{}, fmt.Errorf("incorrect password")
	}
	return ua, nil
}

func register(user string, password string) (UserActivity, error) {
	filePath := "./data/" + user + ".json"
	if _, err := os.Stat(filePath); err == nil {
		return UserActivity{}, fmt.Errorf("user already exists")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return UserActivity{}, fmt.Errorf("error creating file")
	}
	defer file.Close()

	ua := UserActivity{
		Balance:      0,
		Password:     password,
		Transactions: []Transaction{},
	}

	if err := json.NewEncoder(file).Encode(ua); err != nil {
		return UserActivity{}, fmt.Errorf("error writing file")
	}
	return ua, nil
}
