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
	go SendLoginEmail(w, r, db, user)
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

	_, errQuery := db.Exec("INSERT INTO users(email,password,created,usertype,no_card)VALUES(?,?,?,?,?,?)", email, password, created, usertype, noCard)

	if errQuery == nil {

		sendResponse(w, 200, "Success")
	} else {
		fmt.Print(errQuery)
		sendResponse(w, 400, "Insert Failed")
	}
	go SendRegisterEmail(w, r, db, email)
}

func Subscription(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id, _, _ := getUserTokenData(r)

	currentTime := time.Now()
	subscriptionDate := currentTime.Format("2006-01-02")

	rows, err := db.Query("SELECT * FROM users WHERE id=?", id)
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

	fmt.Println(user.Subscription)

	if user.Subscription == nil {
		_, errQuery := db.Exec("UPDATE users SET subscription = ? WHERE id = ?", subscriptionDate, id)
		if errQuery != nil {
			sendResponse(w, 200, "Subscribe Success")
		}
	} else {

		t, err := time.Parse("2006-01-02", *user.Subscription)
		if err != nil {
			fmt.Println(err)
		}
		addedDay := t.AddDate(0, 0, 30)
		_, errQuery := db.Exec("UPDATE users SET subscription = ? WHERE id = ?", addedDay, id)
		if errQuery != nil {
			sendResponse(w, 200, "Subscribe Success")
		}
	}
	go SendSubscriptionEmail(w, r, db, user)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	resetProfileToken(w)
	sendResponse(w, 200, "Logout Success")
}
