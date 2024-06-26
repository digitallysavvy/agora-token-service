package cloud_recording_service

import (
	"math/rand"
	"strconv"
)

// generateUID is a helper function to generate unique user id.
func generateUID() string {
	// Generate a random number between 1 and 2^32 - 1 (4294967295)
	// - starts from 1 to avoid 0, which is reserved.
	uid := rand.Intn(4294967294) + 1

	// Convert the integer to a string
	return strconv.Itoa(uid)
}

// Contains checks if a string is present in a slice.
func Contains(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}
