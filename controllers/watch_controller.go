package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var currState int
var playback int
var watched bool
var watching bool
var paused bool
var profileId int
var videoId string
var liked int
var duration int

func StopWatching() {
	if !watching {
		return
	}
	watching = false

	db := connect()
	defer db.Close()

	var query string
	if !watched {
		query = "INSERT INTO history (id_profile, id_video, latest_state, liked, w_date) VALUES (" + strconv.Itoa(profileId) + ", " + videoId + ", " + strconv.Itoa(currState) + ", " + strconv.Itoa(liked) + ", CURRENT_TIMESTAMP)"
	} else {
		query = "UPDATE history SET latest_state = " + strconv.Itoa(currState) + ", liked = " + strconv.Itoa(liked) + ", w_date = CURRENT_TIMESTAMP WHERE id_profile = " + strconv.Itoa(profileId) + " and id_video = " + videoId
	}
	_, errQuery := db.Exec(query)
	if errQuery == nil {
		fmt.Println("Stop Playing, Successfully updated history")
	} else {
		fmt.Println(errQuery)
		fmt.Println("Stop Playing, Failed updated history")
	}
}

func Player(w http.ResponseWriter, r *http.Request) {
	if !watching {
		sendResponse(w, 400, "Not currently playing a video")
		return
	}

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
		return
	}

	var text string
	action := r.Form.Get("action")
	if action == "forward" {
		currState += 10
		text = "Video forwarded by 10 second"
	} else if action == "backward" {
		currState -= 10
		text = "Video backwarded by 10 second"
	} else if action == "pause" {
		if paused {
			paused = false
			text = "Video is now playing"
		} else {
			paused = true
			text = "Video is now paused"
		}
	} else if action == "like" {
		if liked == 1 {
			liked = 0
			text = "You just disliked the video"
		} else {
			liked = 1
			text = "You just liked the video"
		}
	} else {
		text = "Video is now stopped"
		StopWatching()
	}
	fmt.Println(text)
	sendResponse(w, 200, text)
}

func StartVideo(w http.ResponseWriter, r *http.Request, videoId string, profileId string) {
	db := connect()
	defer db.Close()

	query := "SELECT ms.id, ms.judul, ms.released, ms.age_restriction, ms.sinopsis, ms.genre, ms.pemeran, ms.tags, ms.type, ms.liked, v.id, v.judul, v.description, v.duration, v.season, v.episode FROM movies_and_show ms INNER JOIN video v ON v.id_ms = ms.id WHERE v.id = " + videoId

	rows, err := db.Query(query)

	if err != nil {
		log.Println(err)
		watching = false
		sendResponse(w, 400, "Failed to get movies data")
		return
	}

	var movies MoviesAndShow
	var video Video

	for rows.Next() {
		if err := rows.Scan(&movies.Id, &movies.Judul, &movies.Released, &movies.AgeRestriction, &movies.Sinopsis, &movies.Genre, &movies.Pemeran, &movies.Tags, &movies.MSType, &movies.Liked, &video.Id, &video.JudulVideo, &video.Description, &video.Duration, &video.Season, &video.Episode); err != nil {
			watching = false
			sendResponse(w, 400, "Failed to get movies data")
			log.Println(err)
			return
		} else {
			break
		}
	}

	duration = video.Duration
	var msg string
	if movies.MSType == "Movie" {
		msg = "Playing " + movies.Judul
	} else {
		msg = "Playing " + movies.Judul + " season " + strconv.Itoa(video.Season) + " episode " + strconv.Itoa(video.Episode)
	}
	paused = false
	go WatchTimer(msg)

	sendResponse(w, 200, "Now "+msg)
}

func WatchTimer(txt string) {
	for {
		if !paused {
			currState += playback
			fmt.Println(txt + " at " + strconv.Itoa(currState))
			time.Sleep(1 * time.Second)
		}
		if duration <= currState {
			currState = duration
			StopWatching()
			break
		}
		if !watching {
			break
		}
	}
}

func Watch(w http.ResponseWriter, r *http.Request) {
	StopWatching()
	playback = 1

	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		watching = false
		sendResponse(w, 400, "Failed")
		return
	}

	videoId = r.Form.Get("idvideo")

	fmt.Println(videoId)

	profileId, _ = getProfileTokenData(r)
	id := strconv.Itoa(profileId)
	query := "SELECT latest_state, liked FROM history WHERE id_profile = " + id + " and id_video = " + videoId

	rows, err := db.Query(query)

	if err != nil {
		watching = false
		log.Println(err)
		sendResponse(w, 400, "Something went wrong, please try again.")
		return
	}

	var state int
	var like int

	for rows.Next() {
		if err := rows.Scan(&state, &like); err != nil {
			watching = false
			log.Println(err)
			sendResponse(w, 400, "Something went wrong, please try again.")
			return
		}
	}

	currState = state
	liked = like
	if currState == 0 {
		watched = false
	} else {
		watched = true
	}
	watching = true
	StartVideo(w, r, videoId, id)
}
