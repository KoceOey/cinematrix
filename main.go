package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cinematrix/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")

	router.HandleFunc("/loginProfile", controllers.AuthenticateUser(controllers.ProfileLogin, "Member")).Methods("POST")
	router.HandleFunc("/logoutProfile", controllers.AuthenticateUser(controllers.ProfileLogout, "Member")).Methods("POST")

	router.HandleFunc("/profile", controllers.AuthenticateUser(controllers.ShowProfile, "Member")).Methods("GET")

	router.HandleFunc("/browse", controllers.GetMovies).Methods("GET")
	router.HandleFunc("/watch", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.Watch), "Member")).Methods("POST")
	router.HandleFunc("/watch", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.Player), "Member")).Methods("PUT")

	router.HandleFunc("/createProfile", controllers.CreateProfile).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
