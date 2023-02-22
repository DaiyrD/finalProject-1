package data

import (
	"context"
	"database/sql"
	"finalProjectAdvancedP/internal/validator"
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

func ValidateCart(v *validator.Validator, cart struct {
	UserId   int64 `json:"user_id"`
	BookID   int64 `json:"book_id"`
	Quantity int64 `json:"quantity"`
}) {
	v.Check(cart.UserId > 0, "user_id", "must be greater than zero")
	v.Check(cart.BookID > 0, "book_id", "must be greater than zero")
	v.Check(cart.Quantity > 0, "quantity", "must be greater than zero")
	v.Check(cart.Quantity <= 20, "quantity", "can not be greater than twenty")
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
