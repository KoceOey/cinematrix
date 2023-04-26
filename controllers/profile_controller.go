package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func checkAccountProfileAmount(id int) int {
	db := gormConn()
	var profile []Profile

	db.Where("id_user=?", id).Find(&profile)
	return len(profile)
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	id, _, _ := getUserTokenData(r)

	//Check banyak profile pengguna saat ini
	banyakProfile := checkAccountProfileAmount(id)

	if banyakProfile == 5 { // Jika akun sudah memiliki 5 profile maka tidak boleh create profile lagi
		sendResponse(w, 400, "Sudah mencapai limit profile! (5)")
		return
	}

	StopWatching()

	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendResponse(w, 400, "Failed")
		return
	}

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

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id, _, _ := getUserTokenData(r)

	//Check banyak profile pengguna saat ini
	banyakProfile := checkAccountProfileAmount(id)

	if banyakProfile == 0 { // Jika akun belum memiliki profile maka tidak dapat menghapus
		sendResponse(w, 400, "Belum ada profile di akun ini")
		return
	}

	db := gormConn()
	var profile Profile
	vars := mux.Vars(r)

	id_profile := vars["id_profile"]
	result := db.Where("id=?", id_profile).Delete(&profile)

	if result.RowsAffected < 1 {
		sendResponse(w, 400, "Gagal Delete Profile")
	} else {
		fmt.Println(profile.Nama)
		sendResponse(w, 200, "Berhasil delete profile")
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
