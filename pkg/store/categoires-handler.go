package store

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/go-chi/chi"
)

func (app *OnlineStore) GetCategories(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Query().Get("category_name")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 10
	}

	categories, total, err := app.DB.AllCategories(categoryName, page, pageSize)
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	response := struct {
		Categories []*schema.Category `json:"categories"`
		TotalCount int                `json:"total_count"`
		Page       int                `json:"page"`
		PageSize   int                `json:"page_size"`
		TotalPages int                `json:"total_pages"`
	}{
		Categories: categories,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + 9) / 10,
	}
	app.SendResponse(w, http.StatusOK, response)
}

func (app *OnlineStore) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category schema.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Printf("Error decoding category: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	id, err := app.DB.InsertCategory(&category)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
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

func (app *OnlineStore) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var category schema.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Printf("Error decoding category: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	err = app.DB.UpdateCategory(&category)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, nil)
}

func (app *OnlineStore) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	err = app.DB.DeleteCategory(id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, nil)
}
