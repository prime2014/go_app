package contextkeys

import (
	"errors"
	"net/http"
)

// Define a private custom type to completely prevent collisions
type contextKey string

// Export the key constant
const UserIDKey contextKey = "userID"

func GetAuthenticatedUserID(r *http.Request) (uint, error) {
	userID, ok := r.Context().Value(UserIDKey).(uint)

	if !ok {
		return 0, errors.New("Authentication failed!")
	}

	return uint(userID), nil
}
