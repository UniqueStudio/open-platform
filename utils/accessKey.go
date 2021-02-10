package utils

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

// User is a AccessKey store data type
type User struct {
	UserID  string
	IsAdmin bool
}

// GenAccessKey is a function to gen AccessKey
func GenAccessKey(UserID string, IsAdmin bool) string {
	user := User{
		UserID:  UserID,
		IsAdmin: IsAdmin,
	}

	rs, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}

	AccessKey, err := Encrypt(string(rs))
	if err != nil {
		log.Fatal(err)
	}
	return AccessKey
}

// LoadAccessKey is a func to load AccessKey
func LoadAccessKey(AccessKey string) (userID string, isAdmin bool, err error) {
	if AccessKey == "" {
		return "", false, errors.New("Empyt AccessKey")
	}
	u := new(User)
	data, err := Decrypt(AccessKey)
	if err != nil {
		return "", false, err
	}
	if err := json.Unmarshal([]byte(data), &u); err != nil {
		return "", false, err
	}
	return u.UserID, u.IsAdmin, nil
}
