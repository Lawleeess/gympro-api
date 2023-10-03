package utils

import "strconv"

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
