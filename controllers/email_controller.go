package controllers

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/gomail.v2"
)

func SendLoginEmail(w http.ResponseWriter, r *http.Request, db *sql.DB, user User) {
	mail := gomail.NewMessage()

	mail.SetHeader("From", "cinematrixx@outlook.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "A New Log In")
	htmlBytes, err := ioutil.ReadFile("html/loginSuccess.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	html := string(htmlBytes)
	htmlWithUsername := strings.ReplaceAll(html, "[Email]", user.Email)
	mail.SetBody("text/html", htmlWithUsername)

	dialer := gomail.NewDialer("smtp-mail.outlook.com", 587, "cinematrixx@outlook.com", "Aw1kW0k!!")
	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func SendRegisterEmail(w http.ResponseWriter, r *http.Request, db *sql.DB, email string) {
	mail := gomail.NewMessage()

	mail.SetHeader("From", "cinematrixx@outlook.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "A New Register")
	htmlBytes, err := ioutil.ReadFile("html/registerSuccess.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	html := string(htmlBytes)
	htmlWithUsername := strings.ReplaceAll(html, "[Email]", email)
	mail.SetBody("text/html", htmlWithUsername)

	dialer := gomail.NewDialer("smtp-mail.outlook.com", 587, "cinematrixx@outlook.com", "Aw1kW0k!!")
	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func SendSubscriptionEmail(w http.ResponseWriter, r *http.Request, db *sql.DB, user User) {
	mail := gomail.NewMessage()

	mail.SetHeader("From", "cinematrixx@outlook.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "A New Subscription")
	htmlBytes, err := ioutil.ReadFile("html/subscriptionSuccess.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	html := string(htmlBytes)
	htmlWithUsername := strings.ReplaceAll(html, "[Email]", user.Email)
	htmlWithUsername = strings.ReplaceAll(htmlWithUsername, "[Date]", *user.Subscription)
	mail.SetBody("text/html", htmlWithUsername)

	dialer := gomail.NewDialer("smtp-mail.outlook.com", 587, "cinematrixx@outlook.com", "Aw1kW0k!!")
	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func SendAlmostRanOutEmail(user User) {
	mail := gomail.NewMessage()

	mail.SetHeader("From", "cinematrixx@outlook.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Subscription Almost Ran Out")
	htmlBytes, err := ioutil.ReadFile("html/subscriptionRanOut.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	html := string(htmlBytes)
	htmlWithUsername := strings.ReplaceAll(html, "[Email]", user.Email)
	htmlWithUsername = strings.ReplaceAll(htmlWithUsername, "[Date]", *user.Subscription)
	mail.SetBody("text/html", htmlWithUsername)

	dialer := gomail.NewDialer("smtp-mail.outlook.com", 587, "cinematrixx@outlook.com", "Aw1kW0k!!")
	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func SendSubsRanOutEmail(user User) {
	mail := gomail.NewMessage()

	mail.SetHeader("From", "cinematrixx@outlook.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Subscription Ran Out")
	htmlBytes, err := ioutil.ReadFile("html/subscriptionRanOut.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	html := string(htmlBytes)
	htmlWithUsername := strings.ReplaceAll(html, "[Email]", user.Email)
	htmlWithUsername = strings.ReplaceAll(htmlWithUsername, "[Date]", *user.Subscription)
	htmlWithUsername = strings.ReplaceAll(htmlWithUsername, " almost ", "")
	mail.SetBody("text/html", htmlWithUsername)

	dialer := gomail.NewDialer("smtp-mail.outlook.com", 587, "cinematrixx@outlook.com", "Aw1kW0k!!")
	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
