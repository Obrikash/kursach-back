package main

import "net/http"

func (app *application) listTrainersHandler(w http.ResponseWriter, r *http.Request) {
	trainers, err := app.models.Users.GetTrainers()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"trainers": trainers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
