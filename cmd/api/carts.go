package main

import (
	"errors"
	"finalProjectAdvancedP/internal/data"
	"finalProjectAdvancedP/internal/validator"
	"fmt"
	"net/http"
)

var TotalQuantity int64 = 0

func (app *application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		BookID   int64  `json:"book_id"`
		Quantity int64  `json:"quantity"`
	}

	err := app.readJSON(w, r, &input)
	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if user.Activated != true {
		app.errorResponse(w, r, 404, "Email is not activated or does not exist")
		return
	}

	if user.Activated != true {
		app.errorResponse(w, r, 404, "Email is not activated or does not exist")
		return
	}

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
		Email:         input.Email,
		Quantity:      input.Quantity,
		BookId:        book.ID,
		TotalPrice:    uint64(int64(book.Price) * input.Quantity),
		Books:         books,
		TotalQuantity: TotalQuantity + input.Quantity,
	}

	TotalQuantity = cart.TotalQuantity

	err = app.models.Carts.Insert(cart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cart/%d", cart.Email))
	err = app.writeJSON(w, http.StatusCreated, envelope{"cart": cart}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBookFromCartHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email  string `json:"email"`
		BookID int64  `json:"book_id"`
		Id     int64  `json:"id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		// use our custom error response
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Activated != true {
		app.errorResponse(w, r, 404, "Email is not activated or does not exist")
		return
	}

	cart := &data.Cart{
		Email:  user.Email,
		BookId: input.BookID,
		ID:     input.Id,
	}

	err = app.models.Carts.Delete(cart)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted from the Cart"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listBooksInCartHandler(w http.ResponseWriter, r *http.Request) {

	// Accept the metadata struct as a return value.
	books, err := app.models.Carts.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) orderBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email  string `json:"email"`
		BookID int64  `json:"book_id"`
		Id     int64  `json:"id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	cart := &data.Cart{
		Email:  user.Email,
		BookId: input.BookID,
		ID:     input.Id,
	}

	err = app.models.Carts.Order(cart)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book was successfully ordered"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
