package controllers

type User struct {
	Id           int     `json:"id_user"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	Created      string  `json:"date_created"`
	Subscription *string `json:"subscription_until"`
	UserType     string  `json:"user_type"`
	NoCard       string  `json:"card_number"`
}

type Profile struct {
	Id          int    `json:"id_profile"`
	Nama        string `json:"profile_name"`
	Pin         string `json:"pin"`
	Preferences string `json:"preferences"`
}

type MoviesAndShow struct {
	Id             int         `json:"id_movies_and_show"`
	Judul          string      `json:"judul_movies_and_show"`
	Released       string      `json:"release_date"`
	AgeRestriction int         `json:"age_restriction"`
	Sinopsis       string      `json:"sinopsis"`
	Genre          string      `json:"genre"`
	Pemeran        string      `json:"pemeran"`
	Tags           string      `json:"tags"`
	MSType         string      `json:"type"`
	Liked          int         `json:"type"`
	Videos         interface{} `json:"videos"`
}

type Video struct {
	Id          int    `json:"id_video"`
	JudulVideo  string `json:"judul_video"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Season      int    `json:"season"`
	Episode     int    `json:"episode"`
}

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type MoviesOutput struct {
	Section    string      `json:"section"`
	ListMovies interface{} `json:"listMovies"`
}
