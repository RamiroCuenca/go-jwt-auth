package utils

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generate a random string that contains n amount of digits
func GenerateRandomString(n int) string {

	var str = make([]byte, n)

	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}

	return string(str)
}
