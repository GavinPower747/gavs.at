package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
)

// Basic hasher for generating password hashes for the API_BASIC_AUTH_PASSWORD_HASH environment variable
func main() {
	var password string
	var help bool

	flag.StringVar(&password, "password", "", "Cleartext password to hash")
	flag.StringVar(&password, "p", "", "Cleartext password to hash")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help")

	flag.Parse()

	if help {
		flag.PrintDefaults()

		return
	}

	if password == "" {
		fmt.Println("No password provided")

		return
	}

	fmt.Println(getPasswordHash(password))
}

func getPasswordHash(password string) string {
	hasher := sha256.New()

	hasher.Write([]byte(password))
	bytes := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(bytes)
}
