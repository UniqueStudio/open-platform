package utils

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

// User is a token store data type
type User struct {
	UserID  string
	IsAdmin bool
}

// GenToken is a function to gen token
func GenToken(UserID string, IsAdmin bool) string {
	user := User{
		UserID:  UserID,
		IsAdmin: IsAdmin,
	}

	rs, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}

	Token, err := Encrypt(string(rs))
	if err != nil {
		log.Fatal(err)
	}
	return Token
}

// LoadToken is a func to load token
func LoadToken(Token string) (userID string, isAdmin bool, err error) {
	if Token == "" {
		return "", false, errors.New("Empyt Token")
	}
	u := new(User)
	data, err := Decrypt(Token)
	if err != nil {
		return "", false, err
	}
	if err := json.Unmarshal([]byte(data), &u); err != nil {
		return "", false, err
	}
	return u.UserID, u.IsAdmin, nil
}
