package main

import (
	"fmt"
	"net/http"

	"bookloop.net/internal/data"
	"bookloop.net/internal/validator"
)

func (app *application) createAdvertismentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BookID      int64  `json:"book_id"`
		UserID      int64  `json:"user_id"`
		Publisher   string `json:"publisher"`
		State       string `json:"state"`
		Description string `json:"description,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	advertisment := &data.Advertisment{
		BookID:      input.BookID,
		UserID:      input.UserID,
		Publisher:   input.Publisher,
		State:       input.State,
		Description: input.Description,
	}

	v := validator.New()

	if data.ValidateAdvertisment(v, advertisment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
