package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-crud-api/config"
	"go-crud-api/structures"
	"log"
	"os"
	"time"
)

var DSN = fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%s",
	config.PostgresUser,
	config.PostgresPassword,
	config.DbHost,
	config.DbPort,
	config.PostgresDb)

var Pool *pgxpool.Pool

// Init initializes the database connection
func Init() {
	var err error
	Pool, err = pgxpool.New(context.Background(), DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	log.Println("Connected to the database")
}

// Close closes the database connection
func Close() {
	Pool.Close()
	log.Println("Database connection is closed")
}

// AddItem adds a new item to the database
func AddItem(item structures.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO items (id, name, price) VALUES ($1, $2, $3)"
	_, err := Pool.Exec(ctx, query, item.ID, item.Name, item.Price)
	return err
}

// GetItems retrieves all items from the database
func GetItems() ([]structures.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, price FROM items"
	rows, err := Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []structures.Item
	for rows.Next() {
		var item structures.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetItem retrieves an item by ID
func GetItem(id string) (*structures.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, price FROM items WHERE id=$1"
	row := Pool.QueryRow(ctx, query, id)

	var item structures.Item
	if err := row.Scan(&item.ID, &item.Name, &item.Price); err != nil {
		return nil, err
	}

	return &item, nil
}

// DeleteItem deletes an item by ID
func DeleteItem(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM items WHERE id=$1"
	_, err := Pool.Exec(ctx, query, id)
	return err
}
