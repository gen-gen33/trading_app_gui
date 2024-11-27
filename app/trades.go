package app

import (
	"database/sql"
	"log"
	"trading_app_gui/app/models"
)

// FetchTrades retrieves trade history from the database
func FetchTrades(db *sql.DB) ([]models.Trade, error) {
	query := `
		SELECT t.id, t.buy_order_id, b.user_id AS buyer, t.sell_order_id, s.user_id AS seller, t.amount, t.price, t.created_at
		FROM trades t
		JOIN orders b ON t.buy_order_id = b.id
		JOIN orders s ON t.sell_order_id = s.id
		ORDER BY t.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching trades: %v", err)
		return nil, err
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var trade models.Trade
		err := rows.Scan(
			&trade.ID,
			&trade.BuyOrderID,
			&trade.Buyer,
			&trade.SellOrderID,
			&trade.Seller,
			&trade.Amount,
			&trade.Price,
			&trade.CreatedAt)
		if err != nil {
			log.Printf("Error scanning trade: %v", err)
			continue
		}
		trades = append(trades, trade)
	}
	return trades, nil
}
