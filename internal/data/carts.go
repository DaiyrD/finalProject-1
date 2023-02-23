package data

import (
	"context"
	"database/sql"
	//"finalProjectAdvancedP/internal/validator"
	"github.com/lib/pq"
	"time"
)

type Cart struct {
	ID            int64    `json:"id"`
	Email         string   `json:"email"`
	BookId        int64    `json:"book_id"`
	TotalPrice    uint64   `json:"total_price"`
	Books         []string `json:"books"`
	Quantity      int64    `json:"quantity"`
	TotalQuantity int64    `json:"total_quantity"`
	Ordered       bool     `json:"ordered"`
}

//func ValidateCart(v *validator.Validator, bookID, quantity int64){
//	//UserId   int64 `json:"user_id"`
//	BookID   int64 `json:"book_id"`
//	Quantity int64 `json:"quantity"`
//}) {
//	//v.Check(cart.UserId > 0, "user_id", "must be greater than zero")
//	v.Check(cart.BookID > 0, "book_id", "must be greater than zero")
//	v.Check(cart.Quantity > 0, "quantity", "must be greater than zero")
//	v.Check(cart.Quantity <= 20, "quantity", "can not be greater than twenty")
//}

type CartModel struct {
	DB *sql.DB
}

func (m CartModel) Insert(cart *Cart) error {
	query := `
      INSERT INTO carts (email ,book_id, books, quantity, total_quantity, total_price)
      VALUES ($1, $2, $3, $4, $5, $6)`
	args := []any{cart.Email, cart.BookId, pq.Array(cart.Books), cart.Quantity, cart.TotalQuantity, cart.TotalPrice}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m CartModel) Delete(cart *Cart) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.

	// Construct the SQL query to delete the record.
	query := `
		DELETE from carts
		WHERE email = $1 AND book_id = $2 AND id = $3
`

	args := []any{cart.Email, cart.BookId, cart.ID}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	// Use ExecContext() and pass the context as the first argument.

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m CartModel) GetAll() ([]*Book, error) {
	// Construct the SQL query to retrieve all movie records.
	// Update the SQL query to include the filter conditions.
	// Use full-text search for the title filter.
	// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
	// parameter values.
	query := `
		SELECT * FROM books
		WHERE id IN (SELECT book_id FROM carts)
		ORDER BY title`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	books := []*Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Year,
			&book.Author,
			pq.Array(&book.Genres),
			&book.Price,
			&book.Version,
		)
		if err != nil {
			return nil, err // Update this to return an empty Metadata struct.
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		return nil, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	//metadata := calculateMetadata(totalRecords, f)
	// Include the metadata struct when returning.
	return books, nil
}

func (m CartModel) Order(cart *Cart) error {
	query := `UPDATE carts
	SET ordered = true
	where id = $1 AND email = $2 AND book_id = $3`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{cart.ID, cart.Email, cart.BookId}

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
