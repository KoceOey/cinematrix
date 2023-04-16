package controllers

import (
	"log"
	"net/http"
	"strconv"
)

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id, _, _ := getUserTokenData(r)

	query := "SELECT p.id, p.nama, p.pin, pr.preference FROM profiles p INNER JOIN preferences pr ON p.id = pr.id_profile WHERE id_user = " + strconv.Itoa(id)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}

	var temp Profile
	var profile Profile
	var profiles []Profile
	for rows.Next() {
		if err := rows.Scan(&profile.Id, &profile.Nama, &temp.Pin, &temp.Preferences); err != nil {
			log.Println(err)
			return
		} else {
			profiles = append(profiles, profile)
		}
	}

	sendDataResponse(w, 200, "Success get profile", profiles)
}

// login profile
func ProfileLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
		return
	}

	nama := r.Form.Get("nama")
	pin := r.Form.Get("pin")

	id, _, _ := getUserTokenData(r)

	rows, err := db.Query("SELECT p.id, p.nama, p.pin, pr.preference FROM profiles p INNER JOIN preferences pr ON p.id = pr.id_profile WHERE nama = ? and id_user = ?", nama, id)

	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}

	var profile Profile
	for rows.Next() {
		if err := rows.Scan(&profile.Id, &profile.Nama, &profile.Pin, &profile.Preferences); err != nil {
			log.Println(err)
			return
		} else {
			break
		}
	}

	if pin != profile.Pin {
		sendResponse(w, 400, "Wrong Email/Password!!")
		return
	}

	generateProfileToken(w, profile.Id, profile.Nama)

	// temp send response
	sendResponse(w, 200, "Success login profile")
}

func ProfileLogout(w http.ResponseWriter, r *http.Request) {
	resetProfileToken(w)
	sendResponse(w, 200, "Logout Success")
}
