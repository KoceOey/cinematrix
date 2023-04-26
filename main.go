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

	// endpoint login
	router.HandleFunc("/login", controllers.UserLogin).Methods("POST")

	// endpoint logout
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// endpoint register
	router.HandleFunc("/register", controllers.Register).Methods("POST")

	// endpoint subscription
	router.HandleFunc("/subscription", controllers.AuthenticateUser(controllers.Subscription, "Member")).Methods("POST")

	// endpoint editUser
	router.HandleFunc("/editUser", controllers.AuthenticateUser(controllers.EditUser, "Member")).Methods("PUT")

	//endpoint login profile
	router.HandleFunc("/loginProfile", controllers.AuthenticateUser(controllers.ProfileLogin, "Member")).Methods("POST")

	// endpoint logout profile
	router.HandleFunc("/logoutProfile", controllers.AuthenticateUser(controllers.ProfileLogout, "Member")).Methods("POST")

	// endpoint lihat profile
	router.HandleFunc("/profile", controllers.AuthenticateUser(controllers.ShowProfile, "Member")).Methods("GET")

	// endpoint create profile
	router.HandleFunc("/createProfile", controllers.AuthenticateUser(controllers.CreateProfile, "Member")).Methods("POST")

	// endpoint edit profile
	router.HandleFunc("/editProfile", controllers.AuthenticateUser(controllers.EditProfile, "Member")).Methods("PUT")

	// endpoint delete profile
	router.HandleFunc("/deleteProfile/{id_profile}", controllers.AuthenticateUser(controllers.DeleteProfile, "Member")).Methods("DELETE")

	// endpoint browse / home
	router.HandleFunc("/browse", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.GetMovies), "Member")).Methods("GET")

	// endpoint search
	router.HandleFunc("/search", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.SearchMovie), "Member")).Methods("GET")

	// endpoint watch
	router.HandleFunc("/watch", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.Watch), "Member")).Methods("POST")

	// endpoint player
	router.HandleFunc("/watch", controllers.AuthenticateUser(controllers.AuthenticateProfile(controllers.Player), "Member")).Methods("PUT")

	// endpoint add movie (admin)
	router.HandleFunc("/addMovieShow", controllers.AuthenticateUser(controllers.AddMoviesAndShow, "Admin")).Methods("POST")

	// endpoint add video (admin)
	router.HandleFunc("/addVideo", controllers.AuthenticateUser(controllers.AddVideo, "Admin")).Methods("POST")

	// endpoint remove movie (admin)
	router.HandleFunc("/removeFilm/{id_ms}", controllers.AuthenticateUser(controllers.RemoveFilm, "Admin")).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
