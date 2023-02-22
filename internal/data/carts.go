package data

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Cart struct {
	UserId     int64    `json:"user_id"`
	BookId     int64    `json:"book_id"`
	TotalPrice uint64   `json:"total_price"`
	Books      []string `json:"books"`
	Quantity   int64    `json:"quantity"`
}

type CartModel struct {
	DB *sql.DB
}

func (m CartModel) Insert(cart *Cart) error {
	query := `
      INSERT INTO carts (user_id ,book_id, books, quantity, total_price)
      VALUES ($1, $2, $3, $4, $5)`
	args := []any{cart.UserId, cart.BookId, pq.Array(cart.Books), cart.Quantity, cart.TotalPrice}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
