package store

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/go-chi/chi"
)

func (app *OnlineStore) GetReviewsByProductID(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	reviews, err := app.DB.ReviewsByProductID(productID)
	if err != nil {
		log.Printf("Error getting reviews: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	response := struct {
		Reviews []*schema.Review `json:"reviews"`
	}{
		Reviews: reviews,
	}
	app.SendResponse(w, http.StatusOK, response)
}

func (app *OnlineStore) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review schema.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		log.Printf("Error decoding review: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	id, err := app.DB.InsertReview(&review)
	if err != nil {
		log.Printf("Error inserting review: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	response := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}
	app.SendResponse(w, http.StatusCreated, response)
}

func (app *OnlineStore) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error parsing review ID: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	err = app.DB.DeleteReview(id)
	if err != nil {
		log.Printf("Error deleting review: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, nil)
}
