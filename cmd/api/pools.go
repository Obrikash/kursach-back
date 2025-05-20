package main

import (
	"net/http"

	"github.com/obrikash/swimming_pool/internal/data"
)

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

func (app *application) mostProfitPoolHandler(w http.ResponseWriter, r *http.Request) {
	pool, profit, err := app.models.Pools.MaxProfit()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	mostProfit := struct {
		Pool   data.Pool `json:"pool"`
		Profit float64   `json:"profit"`
	}{
		Pool:   *pool,
		Profit: profit,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"pool": mostProfit}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
