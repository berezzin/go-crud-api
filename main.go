package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/swaggo/http-swagger" // Swagger UI
	_ "go-crud-api/docs"             // Импортируйте документацию
)

// @title CRUD Go API
// @version 1.0
// @description CRUD API with default library
// @host localhost:8080
// @BasePath /

// Item представляет объект данных
type Item struct {
	ID    string  `json:"id"`    // Уникальный идентификатор
	Name  string  `json:"name"`  // Название
	Price float64 `json:"price"` // Цена
}

var (
	items = make(map[string]Item)
	mu    sync.Mutex
)

// getItems возвращает список всех объектов.
// @Summary Возвращает список объектов
// @Tags Items
// @Produce json
// @Success 200 {object} map[string]Item
// @Router /items [get]
func getItems(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(items)
}

// createItem создает новый объект.
// @Summary Создает новый объект
// @Tags Items
// @Accept json
// @Produce json
// @Param item body Item true "Новый объект"
// @Success 201 {object} Item
// @Router /items [post]
func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	items[item.ID] = item
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func main() {
	// Маршруты API
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getItems(w, r)
		} else if r.Method == http.MethodPost {
			createItem(w, r)
		}
	})

	// Swagger-документация
	http.HandleFunc("/docs/", httpSwagger.WrapHandler)

	fmt.Printf("Starting server at port 8080\n")
	http.ListenAndServe(":8080", nil)
}
