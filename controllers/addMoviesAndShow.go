package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func AddMoviesAndShow(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "failed")
		return
	}
	judul := r.Form.Get("judul")
	released := r.Form.Get("released")
	age_restriction := r.Form.Get("age_restriction")
	sinopsis := r.Form.Get("sinopsis")
	genre, _ := strconv.Atoi(r.Form.Get("genre"))
	pemeran := r.Form.Get("pemeran")
	tags := r.Form.Get("tags")
	MSType := r.Form.Get("type")
	Liked := r.Form.Get("liked")

	_, errQuery := db.Exec("INSERT INTO movies_and_show(judul,genre,released,age_restriction,sinopsis,pemeran,tags,type,liked)values (?,?,?,?,?,?,?,?,?)",
		judul,
		genre,
		released,
		age_restriction,
		sinopsis,
		pemeran,
		tags,
		MSType,
		Liked,
	)

	if errQuery == nil {
		sendResponse(w, 200, "Movie added successfully!!!")
	} else {
		sendResponse(w, 400, "Movie added Failed!!!")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sendResponse)
}
