package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	StopWatching()
	db := connect()
	defer db.Close()

	//rediss : "horror, america, anime"
	//preference : [horror america anime]

	var profile Profile
	err := json.Unmarshal([]byte(GetRedis()), &profile)
	if err != nil {
		log.Println(http.StatusBadRequest, "error unmarshal redis")
	}
	delimiter := ", "
	preference := strings.Split(profile.Preferences, delimiter)
	//value : horror, amerika

	// for you
	query := "SELECT ms.id, ms.judul, ms.released, ms.age_restriction, ms.sinopsis, ms.genre, ms.pemeran, ms.tags, ms.type, ms.liked FROM movies_and_show ms WHERE "
	for i, value := range preference {
		if i != 0 {
			query += "or "
		}
		query += "ms.genre LIKE '%" + value + "%' "
	}
	query += "ORDER BY RAND() LIMIT 5"

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

	var listMoviesAndShow []MoviesOutput
	var mnsOutput MoviesOutput
	mnsOutput.Section = "For You"
	mnsOutput.ListMovies = forYous
	listMoviesAndShow = append(listMoviesAndShow, mnsOutput)
	mnsOutput.Section = "Continue Watching"
	mnsOutput.ListMovies = lanjutkanMenontons
	listMoviesAndShow = append(listMoviesAndShow, mnsOutput)
	mnsOutput.Section = "New Release"
	mnsOutput.ListMovies = newReleases
	listMoviesAndShow = append(listMoviesAndShow, mnsOutput)

	sendDataResponse(w, 200, "Success get movies", listMoviesAndShow)
}
