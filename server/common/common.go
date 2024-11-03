package common

import (
	"crypto/rand"
	"fmt"
)

func RandomString(length int) string {
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	return fmt.Sprintf("%x", randomBytes)[:length]
}

type MediaFormat int
