package utils

import (
	"encoding/base64"
	"log"
)

// B64Decode is a func to decode B64 string
func B64Decode(raw string) (string, err) {
	decodeBytes, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(decodeBytes), nil
}

// B64Encode is a func to encode B64 string
func B64Encode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}
