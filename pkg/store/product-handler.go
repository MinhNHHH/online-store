package store

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/go-chi/chi"
)

func (app *OnlineStore) GetProducts(w http.ResponseWriter, r *http.Request) {
	productName := r.URL.Query().Get("product_name")
	categoryName := r.URL.Query().Get("category_name")
	status := r.URL.Query().Get("status")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	products, total, err := app.DB.AllProducts(productName, categoryName, status, page, pageSize)
	if err != nil {
		log.Printf("Error getting products: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	response := struct {
		Products   []*schema.Product `json:"products"`
		TotalCount int               `json:"total_count"`
		Page       int               `json:"page"`
		PageSize   int               `json:"page_size"`
		TotalPages int               `json:"total_pages"`
	}{
		Products:   products,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + 9) / 10,
	}
	app.SendResponse(w, http.StatusOK, response)
}

func (app *OnlineStore) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product schema.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("Error decoding product: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	id, err := app.DB.InsertProduct(&product)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
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

func (app *OnlineStore) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product schema.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	err = app.DB.UpdateProduct(&product)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, nil)
}

func (app *OnlineStore) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	err = app.DB.DeleteProduct(id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		app.SendResponse(w, http.StatusInternalServerError, err)
		return
	}
	app.SendResponse(w, http.StatusOK, nil)
}
