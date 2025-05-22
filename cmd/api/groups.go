package main

import (
	"fmt"
	"net/http"

	"github.com/obrikash/swimming_pool/internal/data"
)

func (app *application) listGroupsHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := app.models.Groups.GetGroups()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"groups": groups}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addGroupToPoolHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PoolID     int64 `json:"pool_id"`
		CategoryID int64 `json:"category_id"`
		TrainerID  int64 `json:"trainer_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	group := data.Group{Pool: input.PoolID, Category: input.CategoryID, Trainer: data.User{ID: input.TrainerID}}

	err = app.models.Groups.AddToPool(group)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"success": fmt.Sprintf("group is created with ID %d", group.ID)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
