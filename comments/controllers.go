package comments

import (
	contextkeys "contextKeys"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentController struct {
	Service CommentService
}

func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(uint)

	if !ok || userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized! You are not authorized to write comments"})
		return
	}

	vars := mux.Vars(r)

	blogIdStr := vars["blogId"]

	blogId, err := strconv.ParseUint(blogIdStr, 10, 32)

	if err != nil || blogId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or missing Blog ID path parameter"})
		return
	}

	blogID := uint(blogId)

	var dto CommentDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not parse the comment dto"})
		return
	}

	comment, err := c.Service.CreateComment(dto, blogID, userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)

}
