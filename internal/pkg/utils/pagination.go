package utils

import (
	"crypto/rand"
	"io"
	"strconv"
)

const (
	defaultOffset = 0
	defaultLimit  = 20
)

// IsValidPage verifies if the given offset and limit are not empty and are valid numbers.
func IsValidPage(offset, limit string) bool {
	if offset == "" || limit == "" {
		return false
	}

	o, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return false
	}

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return false
	}

	return l > 0 && o >= 0
}

// GetValidPage returns default offset (0) and default limit (100) to paginate results.
func GetValidPage() (int64, int64) {
	return defaultOffset, defaultLimit
}

func EncodeToString(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
