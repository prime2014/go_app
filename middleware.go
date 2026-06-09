package main

import (
	"context"
	contextkeys "contextKeys"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Define your JWT claims signature
type MyCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Log the incoming request details
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		//2. Pass control to the next handler in the chain
		next.ServeHTTP(w, r)

		//3. Log the completion and duration
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// 1. Check if the Authorization header was provided at all
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Parse out the "Bearer " prefix safely
		if !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "Unauthorized: Authentication must use bearer token format", http.StatusUnauthorized)
		}

		// Extract the actual token string by stripping out the first 7 characters
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			http.Error(w, "Unauthorized: Token payload cannot be empty", http.StatusUnauthorized)
			return
		}

		// 3. Validate your token (Placeholder logic for your JWT or session lookup)
		userID, err := validateYourTokenString(tokenString)

		if err != nil {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		//4. OPTIONAL: Inject your metadata into the request context
		// This makes the logged in User ID available to any route handlers
		ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)

		// 5. Pass control to the next handler using the enrichec context
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func validateYourTokenString(tokenString string) (uint, error) {
	err := godotenv.Load()

	if err != nil {
		return 0, err
	}

	signingKey := os.Getenv("SECRET_KEY")

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	// Extract claims and validate validity
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {

		return claims.UserID, nil
	}

	return 0, fmt.Errorf("invalid token")

}
