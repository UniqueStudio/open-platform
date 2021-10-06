package utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// State is resp state type
type State struct {
	URL       string `json:"url"`
	Token     string `json:"token"`
	AccessKey string `json:"accessKey"`
}

// GenToken is a func to generate token
func GenToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"expire": func() int64 {
			now := time.Now()
			duration, _ := time.ParseDuration("14d")
			m1 := now.Add(duration)
			return m1.Unix()
		}(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(AppConfig.Server.SecretKey))

	return tokenString, err
}
