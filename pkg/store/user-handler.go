package store

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (app *OnlineStore) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	userID := request.UserID
	productID := request.ProductID
	err := app.DB.AddToWishlist(userID, productID)
	if err != nil {
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusCreated, nil)
}

func (app *OnlineStore) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	userID := request.UserID
	productID := request.ProductID
	err := app.DB.RemoveFromWishlist(userID, productID)
	if err != nil {
		log.Printf("Error removing from wishlist: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, nil)
}

func (app *OnlineStore) GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	products, err := app.DB.GetWishlist(userID)
	if err != nil {
		log.Printf("Error getting wishlist: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, products)
}
