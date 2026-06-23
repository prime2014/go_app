package blogs

import (
	contextkeys "contextKeys"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

type BlogController struct {
	Service BlogServices
}

func (b *BlogController) CreateBlog(w http.ResponseWriter, r *http.Request) {

	userID, _ := contextkeys.GetAuthenticatedUserID(r)

	if userID == 0 {
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

func (b *BlogController) EditBlog(w http.ResponseWriter, r *http.Request) {
	userID, _ := contextkeys.GetAuthenticatedUserID(r)

	if userID == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "This Action is Unauthorized!"})
	}

	var editBlogDto EditBlogDto
	vars := mux.Vars(r)

	blogID, _ := strconv.Atoi(vars["blogID"])

	if err := json.NewDecoder(r.Body).Decode(&editBlogDto); err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "The payload you entered cannot be processed!"})
	}

	if err := validate.Struct(&editBlogDto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation error: " + err.Error()})
	}

	blog, err := b.Service.Edit(editBlogDto, userID, uint(blogID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error: " + err.Error()})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&blog)

}
