package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"priceCents"`
	Stock       int    `json:"stock"`
}

type Cart struct {
	SessionID string `json:"sessionid" gorm:"primaryKey"`
	ProductID uint   `json:"productid" gorm:"primaryKey"`
	Quantity  int    `json:"quantity"`
}

type CartResponse struct {
	SessionID   string `json:"sessionId"`
	ProductID   uint   `json:"productId"`
	ProductName string `json:"productName"`
	PriceCents  int    `json:"priceCents"`
	Quantity    int    `json:"quantity"`
}

func main() {

	db, err := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{}, &Cart{})
	var count int64
	db.Model(&Product{}).Count(&count)
	if count == 0 {
		db.Create(&Product{Name: "Product 1", Description: "Description 1", PriceCents: 1000, Stock: 10})
		db.Create(&Product{Name: "Product 2", Description: "Description 2", PriceCents: 2000, Stock: 20})
	}
	log.Info().Msg("Data seeded")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		// handler: return products as JSON
		var products []Product
		result := db.Find(&products)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	r.Get("/cart/{sessionId}", func(w http.ResponseWriter, r *http.Request) {
		// handler: return cart

		sessionId := chi.URLParam(r, "sessionId")
		var cart []CartResponse
		result := db.Table("carts").
			Select("carts.session_id, carts.product_id, products.name as product_name, products.price_cents, carts.quantity").
			Joins("JOIN products ON carts.product_id = products.id").
			Where("carts.session_id = ?", sessionId).Scan(&cart)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cart)
	})

	r.Post("/cart/items", func(w http.ResponseWriter, r *http.Request) {
		// handler: add item to cart
		var item Cart
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingItem Cart
		result := db.Where("session_id = ? AND product_id = ?", item.SessionID, item.ProductID).First(&existingItem)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			db.Create(&item)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(item)
			return
		}
		if result.Error == nil {
			existingItem.Quantity += item.Quantity
			db.Save(&existingItem)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingItem)
			return
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Delete("/cart/{sessionId}", func(w http.ResponseWriter, r *http.Request) {
		sessionId := chi.URLParam(r, "sessionId")

		// Delete all cart rows for this session
		result := db.Where("session_id = ?", sessionId).Delete(&Cart{})
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "Cart cleared",
			"deletedRows": result.RowsAffected,
		})
	})

	http.ListenAndServe(":8080", r)

}
