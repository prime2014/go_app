package contextkeys

// Define a private custom type to completely prevent collisions
type contextKey string

// Export the key constant
const UserIDKey contextKey = "userID"
