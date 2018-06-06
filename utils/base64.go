package utils

import (
	"encoding/base64"
	"log"
)

// B64Decode is a func to decode B64 string
func B64Decode(raw string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		log.Fatalln(err)
	}
	return string(decodeBytes)
}
