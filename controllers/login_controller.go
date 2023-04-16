package controllers

import (
	"log"
	"net/http"
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

	if password != user.Password{
		sendResponse(w, 400, "Wrong Email/Password!!")
		return
	}

	generateUserToken(w, user.Id, user.Email, user.UserType)
	if(user.UserType == "Member"){
		// show profile
	}else{
		// show menu admin
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request){
	resetUserToken(w)
	sendResponse(w, 200, "Logout Success")
}