package data

import (
	"context"
	"database/sql"
	"errors"
	"finalProjectAdvancedP/internal/validator"
	"fmt"
	"github.com/lib/pq"
	"time"
)

// Annotate the Book struct with struct tags to control how the keys appear in the
// JSON-encoded output.

type Book struct {
	Author    string    `json:"author"`
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive in order to hide this info from JSON response
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`   // Add the omitempty directive to output this info only if it is not empty
	Genres    []string  `json:"genres,omitempty"` // Add the omitempty directive
	Price     uint64    `json:"price"`
	Version   int32     `json:"version"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Author != "", "author", "must be provided")
	v.Check(len(book.Author) <= 100, "author", "must not be more than 100 bytes long")
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year >= 1455, "year", "must be greater than 1455")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
	v.Check(book.Price > 0, "price", "must be greater than zero")
}

// Define a BookModel struct type which wraps a sql.DB connection pool.
type BookModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
// The Insert() method accepts a pointer to a movie struct, which should contain the
// data for the new record.
func (m BookModel) Insert(book *Book) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO books (title, year, genres, author, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.

	// Create a context with a 3-second timeout.
	args := []any{book.Title, book.Year, pq.Array(book.Genres), book.Author, book.Price}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)

}

// Add a placeholder method for fetching a specific record from the movies table.
func (m BookModel) Get(id int64) (*Book, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
		SELECT id, created_at, title, year, author, genres, price, version
		FROM books
		WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var book Book

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.

	// Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Year,
		&book.Author,
		pq.Array(&book.Genres),
		&book.Price,
		&book.Version,
	)

	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &book, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m BookModel) Update(book *Book) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
		UPDATE books
		SET title = $1, year = $2, author = $3, genres = $4, price = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		book.Title,
		book.Year,
		book.Author,
		pq.Array(book.Genres),
		book.Price,
		book.ID,
		book.Version,
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRow() method to execute the query, passing in the args slice as a
	// variadic parameter and scanning the new version value into the movie struct.
	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM books
		WHERE id = $1`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	//result, err := m.DB.Exec(query, id)
	//if err != nil {
	//	return err
	//}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (m BookModel) GetAll(title string, genres []string, filters Filters) ([]*Book, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	// Update the SQL query to include the filter conditions.
	// Use full-text search for the title filter.
	// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
	// parameter values.
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, year, author, genres, price, version
FROM books
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	books := []*Book{}
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
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
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return books, metadata, nil
}
