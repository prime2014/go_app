package accounts

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var validate = validator.New()

type MyCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type UserController struct {
	Service UserService
}

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateToken(userID uint) (string, error) {

	err := godotenv.Load()

	if err != nil {
		return "", err
	}

	signingKey := []byte(os.Getenv("SECRET_KEY"))

	claims := MyCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        generateJTI(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "xpeng-app",
			Audience:  jwt.ClaimStrings{"xpeng-client"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Removed the "error" return type from the signature
func (u *UserController) RegisterUsers(w http.ResponseWriter, r *http.Request) {

	var dto SignUpDto
	log.Println("Inside RegisterUsers controller")
	// Set JSON header for responses
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return // Just return to stop execution
	}

	log.Println(dto)

	if err := validate.Struct(dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed: " + err.Error()})
		return
	}

	user, err := u.Service.SignupUser(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register user: " + err.Error()})
		return
	}

	log.Println(user)

	// Properly encode success as JSON instead of fmt.Fprintln text
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "Account created successfully",
		"email":   user.Email,
	})
}

func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginDto LoginDto

	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to login!"})
		return
	}

	user, err := u.Service.LoginUser(loginDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	tokenString, err := generateToken(user.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "token serialization error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"access_token": tokenString})
}
