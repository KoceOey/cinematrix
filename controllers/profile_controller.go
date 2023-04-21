package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var context_redis = context.Background()

func RedisInit() {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rdb = db
}

func GetRedis() string {

	val, err := rdb.Get(context_redis, "key").Result()
	if err == redis.Nil {
		log.Println(http.StatusNotFound, "data tidak ditemukan")
	} else if err != nil {
		log.Println(http.StatusBadRequest, "error get redis")
	}

	return val
}

func GetRedisTest(key string) string {

	val, err := rdb.Get(context_redis, key).Result()
	if err == redis.Nil {
		log.Println(http.StatusNotFound, "data tidak ditemukan")
	} else if err != nil {
		log.Println(http.StatusBadRequest, "error get redis")
	}

	return val
}

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	StopWatching()

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
		if err := rows.Scan(&profile.Id, &profile.Nama, &temp.Pin, &profile.Preferences); err != nil {
			log.Println(err)
			return
		} else {
			profiles = append(profiles, profile)
		}
	}

	sendDataResponse(w, 200, "Success get profile", profiles)
}

func ShowAdmin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	//rediss : "horror, america, anime"
	//preference : [horror america anime]

	// var profile Profile
	// err := json.Unmarshal([]byte(GetRedis()), &profile)
	// if err != nil {
	// 	log.Println(http.StatusBadRequest, "error unmarshal redis")
	// }
	// delimiter := ", "
	// preference := strings.Split(profile.Preferences, delimiter)
	//value : horror, amerika

	// for you
	query := "SELECT ms.id, ms.judul, ms.released, ms.age_restriction, ms.sinopsis, ms.genre, ms.pemeran, ms.tags, ms.type, ms.liked FROM movies_and_show ms "
	// for i, value := range preference {
	// 	if i != 0 {
	// 		query += "or "
	// 	}
	// 	query += "ms.genre LIKE '%" + value + "%' "
	// }
	// query += "ORDER BY RAND() LIMIT 5"

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}
	var forYou MoviesAndShow
	var forYous []MoviesAndShow
	for rows.Next() {
		if err := rows.Scan(&forYou.Id, &forYou.Judul, &forYou.Released, &forYou.AgeRestriction, &forYou.Sinopsis, &forYou.Genre, &forYou.Pemeran, &forYou.Tags, &forYou.MSType, &forYou.Liked); err != nil {
			log.Println(err)
			return
		} else {
			query2 := "SELECT v.id, v.judul, v.description, v.duration, v.season, v.episode FROM video v WHERE v.id_ms = " + strconv.Itoa(forYou.Id)
			result, errQuery := db.Query(query2)

			if errQuery != nil {
				log.Println(err)
				sendResponse(w, 400, "Something went wrong, please try again.")
				return
			}

			var video Video
			var videos []Video
			for result.Next() {
				if err := result.Scan(&video.Id, &video.JudulVideo, &video.Description, &video.Duration, &video.Season, &video.Episode); err != nil {
					log.Println(err)
					return
				} else {
					videos = append(videos, video)
				}
			}
			forYou.Videos = videos
			forYous = append(forYous, forYou)
		}
	}

	// lanjutkan menonton

	id, _ := getProfileTokenData(r)

	query = "SELECT ms.id, ms.judul, ms.released, ms.age_restriction, ms.sinopsis, ms.genre, ms.pemeran, ms.tags, ms.type, ms.liked, v.id, v.judul, v.description, v.duration, v.season, v.episode FROM movies_and_show ms INNER JOIN video v ON v.id_ms = ms.id INNER JOIN history h ON h.id_video = v.id WHERE h.id_profile = " + strconv.Itoa(id) + " ORDER BY h.w_date DESC LIMIT 5"

	rows, err = db.Query(query)

	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}

	var lanjutkanMenonton MoviesAndShow
	var video Video
	var lanjutkanMenontons []MoviesAndShow

	for rows.Next() {
		if err := rows.Scan(&lanjutkanMenonton.Id, &lanjutkanMenonton.Judul, &lanjutkanMenonton.Released, &lanjutkanMenonton.AgeRestriction, &lanjutkanMenonton.Sinopsis, &lanjutkanMenonton.Genre, &lanjutkanMenonton.Pemeran, &lanjutkanMenonton.Tags, &lanjutkanMenonton.MSType, &lanjutkanMenonton.Liked, &video.Id, &video.JudulVideo, &video.Description, &video.Duration, &video.Season, &video.Episode); err != nil {
			log.Println(err)
			return
		} else {
			lanjutkanMenonton.Videos = video
			lanjutkanMenontons = append(lanjutkanMenontons, lanjutkanMenonton)
		}
	}

	// baru rilis
	query = "SELECT ms.id, ms.judul, ms.released, ms.age_restriction, ms.sinopsis, ms.genre, ms.pemeran, ms.tags, ms.type, ms.liked FROM movies_and_show ms ORDER BY ms.released DESC LIMIT 5"

	rows, err = db.Query(query)

	if err != nil {
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}

	var newRelease MoviesAndShow
	var newReleases []MoviesAndShow

	for rows.Next() {
		if err := rows.Scan(&newRelease.Id, &newRelease.Judul, &newRelease.Released, &newRelease.AgeRestriction, &newRelease.Sinopsis, &newRelease.Genre, &newRelease.Pemeran, &newRelease.Tags, &newRelease.MSType, &newRelease.Liked); err != nil {
			log.Println(err)
			return
		} else {
			query2 := "SELECT v.id, v.judul, v.description, v.duration, v.season, v.episode FROM video v WHERE v.id_ms = " + strconv.Itoa(newRelease.Id)
			result, errQuery := db.Query(query2)

			if errQuery != nil {
				log.Println(err)
				sendResponse(w, 400, "Something went wrong, please try again.")
				return
			}

			var video Video
			var videos []Video
			for result.Next() {
				if err := result.Scan(&video.Id, &video.JudulVideo, &video.Description, &video.Duration, &video.Season, &video.Episode); err != nil {
					log.Println(err)
					return
				} else {
					videos = append(videos, video)
				}
			}
			newRelease.Videos = videos
			newReleases = append(newReleases, newRelease)
		}
	}

	var listMoviesAndShow []ResponseLoginAdmin
	var mnsOutput ResponseLoginAdmin
	mnsOutput.Section = "For You"
	mnsOutput.ListMovies = forYous
	listMoviesAndShow = append(listMoviesAndShow, mnsOutput)
	mnsOutput.Section = "Continue Watching"
	mnsOutput.ListMovies = lanjutkanMenontons
	listMoviesAndShow = append(listMoviesAndShow, mnsOutput)
	mnsOutput.Section = "New Release"
	mnsOutput.ListMovies = newReleases
	listMoviesAndShow = append(listMoviesAndShow, mnsOutput)

	sendDataResponse(w, 200, "Success Login", listMoviesAndShow)
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	StopWatching()

	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
		return
	}

	id, _, _ := getUserTokenData(r)

	nama := r.Form.Get("profile_name")
	pin := r.Form.Get("pin")

	_, errQuery := db.Exec("INSERT INTO profiles(id_user,nama,pin) VALUES (?,?,?)", id, nama, pin)
	if errQuery == nil {
		sendResponse(w, 200, "Success")
	} else {
		fmt.Print(errQuery)
		sendResponse(w, 400, "Insert Failed")
	}
}

// login profile
func ProfileLogin(w http.ResponseWriter, r *http.Request) {
	StopWatching()
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

	reqRedis := Profile{
		Preferences: profile.Preferences,
	}
	req, _ := json.Marshal(reqRedis)
	RedisInit()
	errSet := rdb.Set(context_redis, "key", req, 0).Err()

	if errSet != nil {
		log.Println("Error Set Redis", errSet)
	}

	// temp send response
	GetMovies(w, r)
}

func ProfileLogout(w http.ResponseWriter, r *http.Request) {
	StopWatching()
	resetProfileToken(w)
	sendResponse(w, 200, "Logout Success")
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
	}

	nama := r.Form.Get("nama")
	fmt.Println(nama)
	pin := r.Form.Get("pin")
	fmt.Println(pin)

	id, _, _ := getUserTokenData(r)
	idProfile, namaProfile := getProfileTokenData(r)

	if nama == "" {

		_, errQuery := db.Exec("UPDATE profiles set pin=? WHERE id_user = ? AND id =?", pin, id, idProfile)

		if err != nil {
			log.Println(errQuery)
			sendResponse(w, 400, "Something went wrong, please try again.")
			return
		}
	} else if pin == "" {
		_, errQuery := db.Exec("UPDATE profiles set nama=? WHERE id_user = ? AND id =?", nama, id, idProfile)

		if err != nil {
			log.Println(errQuery)
			sendResponse(w, 400, "Something went wrong, please try again.")
			return
		}
	} else if pin == "" && nama == "" {
		sendDataResponse(w, 400, "No Data Input For edit profile ", namaProfile)

	} else {
		_, errQuery := db.Exec("UPDATE profiles set nama=?, pin=? WHERE id_user = ? AND id =?", nama, pin, id, idProfile)

		if err != nil {
			log.Println(errQuery)
			sendResponse(w, 400, "Something went wrong, please try again.")
			return
		}
	}

	sendDataResponse(w, 200, "Success edit profile", namaProfile)
}
