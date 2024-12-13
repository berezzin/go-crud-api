package structures

// Item data object
type Item struct {
	ID    string  `json:"id"`    // Unique ID
	Name  string  `json:"name"`  // Object's name
	Price float64 `json:"price"` // Price
}
