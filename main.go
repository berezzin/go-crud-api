package main

import (
	"encoding/json"
	"fmt"
	"go-crud-api/db"
	"go-crud-api/structures"
	"log"
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

var (
	items = make(map[string]structures.Item)
	mutex sync.Mutex
)

// getItems Get list of Items.
// @Summary Return list of all Items.
// @Tags Items
// @Produce json
// @Success 200 {object} map[string]structures.Item
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
// @Param item body structures.Item true "New Item"
// @Success 201 {object} structures.Item
// @Router /items [post]
func createItem(w http.ResponseWriter, r *http.Request) {
	var item structures.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	items[item.ID] = item
	mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// getItem get Item by ID.
// @Summary get Item by ID.
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Object's ID"
// @Success 200 {object} structures.Item
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

// deleteItem Delete Item by ID.
// @Summary delete Item by ID.
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Object's ID"
// @Success 200 {object} structures.Response
// @Failure 404 {string} string "Item not found"
// @Failure 400 {string} string "Bad request"
// @Router /items/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/items/"):]
	if id == "" {
		http.Error(w, "ID didn't provided", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	_, exists := items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	delete(items, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := structures.Response{
		Status: "ok",
		Detail: fmt.Sprintf("Item with id '%s' deleted", id),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	db.Init()
	defer db.Close()
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
		} else if r.Method == http.MethodDelete {
			deleteItem(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Swagger-документация
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)

	log.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", nil)
}
