package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/swaggo/http-swagger"
	_ "go-crud-api/docs"
)

// @title CRUD Go API
// @version 1.0
// @description CRUD API with default library
// @host localhost:8080
// @BasePath /

// Item data object
type Item struct {
	ID    string  `json:"id"`    // Unique ID
	Name  string  `json:"name"`  // Object's name
	Price float64 `json:"price"` // Price
}

var (
	items = make(map[string]Item)
	mutex sync.Mutex
)

// getItems Get list of Items.
// @Summary Return list of all Items.
// @Tags Items
// @Produce json
// @Success 200 {object} map[string]Item
// @Router /items [get]
func getItems(w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// createItem Create a new Item.
// @Summary Create a new Item object.
// @Tags Items
// @Accept json
// @Produce json
// @Param item body Item true "New Item"
// @Success 201 {object} Item
// @Router /items [post]
func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	items[item.ID] = item
	mutex.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// getItem get Item by ID.
// @Summary get Item by ID.
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Object's ID"
// @Success 200 {object} Item
// @Failure 404 {string} string "Item not found"
// @Failure 400 {string} string "Bad request"
// @Router /items/{id} [get]
func getItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/items/"):]
	if id == "" {
		http.Error(w, "ID didn't provided", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	item, exists := items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func main() {
	// Маршруты API
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getItems(w)
		} else if r.Method == http.MethodPost {
			createItem(w, r)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getItem(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Swagger-документация
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)

	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)
}
