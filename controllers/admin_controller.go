package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	genre := r.Form.Get("genre")
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

func AddVideo(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "failed")
		return
	}

	id_ms := r.Form.Get("id_ms")
	judul := r.Form.Get("judul")
	description := r.Form.Get("description")
	duration := r.Form.Get("duration")
	season, _ := strconv.Atoi(r.Form.Get("season"))
	episode := r.Form.Get("episode")

	_, errQuery := db.Exec("INSERT INTO video(id_ms,season,judul,description,duration,episode)values (?,?,?,?,?,?)",
		id_ms,
		season,
		judul,
		description,
		duration,
		episode,
	)

	if errQuery == nil {
		sendResponse(w, 200, "Movie added successfully!!!")
	} else {
		sendResponse(w, 400, "Movie added Failed!!!")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sendResponse)

}

func RemoveFilm(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	id_ms := vars["id_ms"]
	fmt.Println(id_ms)

	query := "select id_ms FROM video WHERE id_ms = '" + id_ms + "'"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again")
		return
	}

	var video Video
	for rows.Next() {
		fmt.Println(video.Id_ms)
		if err := rows.Scan(&video.Id_ms); err != nil {
			log.Println(err)
			sendResponse(w, 400, "Something went wrong, please try again")
			return
		} else {
			sendResponse(w, 400, "Data movies and show already on video!!!")
			return
		}
	}

	_, errQuery := db.Exec("DELETE FROM movies_and_show WHERE id=?",
		id_ms,
	)

	if errQuery == nil {
		sendResponse(w, 200, "Movie removed successfully!!!")
	} else {
		sendResponse(w, 400, "Movie remove Failed!!!")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sendResponse)

}
