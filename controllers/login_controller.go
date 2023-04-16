package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	rows, err := db.Query("SELECT * FROM users WHERE email = ?", email)

	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}

	var user User
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Created, &user.Subscription, &user.UserType, &user.NoCard); err != nil {
			log.Println(err)
			return
		} else {
			break
		}
	}

	if password != user.Password {
		sendResponse(w, 400, "Wrong Email/Password!!")
		return
	}

	generateUserToken(w, user.Id, user.Email, user.UserType)
	if user.UserType == "Member" {
		ShowProfile(w, r)
	} else {
		// show menu admin
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	currentTime := time.Now()
	created := currentTime.Format("2006-01-02")

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
	}

	email := r.Form.Get("email")
	fmt.Print("email : ", email)
	password := r.Form.Get("password")
	noCard := r.Form.Get("no_card")
	usertype := "Member"

	_, errQuery := db.Exec("INSERT INTO users(email,password,created,subscription,usertype,no_card)VALUES(?,?,?,?,?,?)", email, password, created, created, usertype, noCard)

	if errQuery == nil {

		sendResponse(w, 200, "Success")
	} else {
		fmt.Print(errQuery)
		sendResponse(w, 400, "Insert Failed")
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	resetProfileToken(w)
	sendResponse(w, 200, "Logout Success")
}
