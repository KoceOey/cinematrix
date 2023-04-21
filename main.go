package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cinematrix/controllers"
	"github.com/go-co-op/gocron"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// CRON
	s := gocron.NewScheduler(time.UTC)
	//Day().At("12:00")
	s.Every(1).Day().Do(controllers.Task)
	s.StartAsync()

	router := mux.NewRouter()

	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/subscription", controllers.AuthenticateUser(controllers.Subscription, "Member")).Methods("POST")
	router.HandleFunc("/editUser", controllers.AuthenticateUser(controllers.EditUser, "Member")).Methods("PUT")

	router.HandleFunc("/loginProfile", controllers.AuthenticateUser(controllers.ProfileLogin, "Member")).Methods("POST")
	router.HandleFunc("/logoutProfile", controllers.AuthenticateUser(controllers.ProfileLogout, "Member")).Methods("POST")
	router.HandleFunc("/profile", controllers.AuthenticateUser(controllers.ShowProfile, "Member")).Methods("GET")
	router.HandleFunc("/createProfile", controllers.AuthenticateUser(controllers.CreateProfile, "Member")).Methods("POST")
	router.HandleFunc("/editProfile", controllers.AuthenticateUser(controllers.EditProfile, "Member")).Methods("PUT")

	router.HandleFunc("/browse", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.GetMovies), "Member")).Methods("GET")
	router.HandleFunc("/search", controllers.SearchMovie).Methods("GET")
	router.HandleFunc("/watch", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.Watch), "Member")).Methods("POST")
	router.HandleFunc("/watch", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.Player), "Member")).Methods("PUT")

	router.HandleFunc("/addMovieShow", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.AddMoviesAndShow), "Admin")).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
