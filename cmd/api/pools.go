package main

import "net/http"

func (app *application) listPoolsHandler(w http.ResponseWriter, r *http.Request) {
	pools, err := app.models.Pools.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"pools": pools}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
