package store

import (
	"encoding/json"
	"net/http"

	"github.com/MinhNHHH/online-store/pkg/cfgs"
	databases "github.com/MinhNHHH/online-store/pkg/databases/repositories"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type OnlineStore struct {
	Cfgs    cfgs.Configs
	DB      databases.DatabaseRepo
	Session *scs.SessionManager
}

func (app *OnlineStore) SendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (app *OnlineStore) Routes() http.Handler {
	mux := chi.NewRouter()
	// register middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	// register routes
	mux.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth", app.authenticate)
		r.Post("/refresh-token", app.refresh)
		r.Post("/auth/register", app.CreateUser)
		r.Route("/users", func(rUser chi.Router) {
			// rUser.Put("/{id}", app.UpdateUser)
			// rUser.Delete("/{id}", app.DeleteUser)
			rUser.Route("/wishlist", func(rWishlist chi.Router) {
				rWishlist.Use(app.authRequired)
				rWishlist.Post("/", app.AddToWishlist)
				rWishlist.Delete("/", app.RemoveFromWishlist)
				rWishlist.Get("/", app.GetWishlist)
			})
		})
		r.Route("/products", func(rProduct chi.Router) {
			rProduct.Use(app.authRequired)
			rProduct.Get("/", app.GetProducts)
			rProduct.Post("/", app.CreateProduct)
			rProduct.Put("/{id}", app.UpdateProduct)
			rProduct.Delete("/{id}", app.DeleteProduct)
		})
		r.Route("/categories", func(rCategory chi.Router) {
			rCategory.Use(app.authRequired)
			rCategory.Get("/", app.GetCategories)
			rCategory.Post("/", app.CreateCategory)
			rCategory.Put("/{id}", app.UpdateCategory)
			rCategory.Delete("/{id}", app.DeleteCategory)
		})
		r.Route("/reviews", func(rReview chi.Router) {
			rReview.Use(app.authRequired)
			rReview.Get("/{product_id}", app.GetReviewsByProductID)
			rReview.Post("/{product_id}", app.CreateReview)
			rReview.Delete("/{product_id}", app.DeleteReview)
		})
	})

	return mux
}
