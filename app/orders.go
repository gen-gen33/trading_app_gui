package app

import (
	"database/sql"
	"log"
	"trading_app_cli/models"
)

// FetchOrders retrieves open orders from the database
func FetchOrders(db *sql.DB) ([]models.Order, error) {
	query := `
		SELECT id, user_id, type, amount, price
		FROM orders
		WHERE status = 'open'
		ORDER BY created_at ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.User, &order.Type, &order.Amount, &order.Price)
		if err != nil {
			log.Printf("Error scanning order: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}
