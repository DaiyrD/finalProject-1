package main

import (
	"errors"
	"finalProjectAdvancedP/internal/data"
	"fmt"
	"net/http"
)

func (app *application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId   int64 `json:"user_id"`
		BookID   int64 `json:"book_id"`
		Quantity int64 `json:"quantity"`
	}
	err := app.readJSON(w, r, &input)
	book, err := app.models.Books.Get(input.BookID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if err != nil { // use our custom error response
		app.badRequestResponse(w, r, err)
		return
	}
	books := make([]string, 0)
	books = append(books, book.Title)
	cart := &data.Cart{
		UserId:     input.UserId,
		Quantity:   input.Quantity,
		BookId:     book.ID,
		TotalPrice: uint64(int64(book.Price) * input.Quantity),
		Books:      books,
	}
	err = app.models.Carts.Insert(cart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cart/%d", cart.UserId))
	err = app.writeJSON(w, http.StatusCreated, envelope{"cart": cart}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
