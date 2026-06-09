package blogs

import (
	contextkeys "contextKeys"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type BlogController struct {
	Service BlogServices
}

func (b *BlogController) CreateBlog(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Context())
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(uint)
	fmt.Println("The userID is: ", userID)

	if !ok || userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized! You are not authorized to create a blog"})
		return
	}

	var dto BlogDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "error in serialization of the request body"})
	}

	if err := validate.Struct(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation error: " + err.Error()})
	}

	blog, err := b.Service.Create(dto, userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&blog)

}
