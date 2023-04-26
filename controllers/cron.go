package controllers

import (
	"fmt"
	"math"
	"time"
)

func Task() {
	users := getAllUsers()
	CheckSub(users)
}
func getAllUsers() []User {
	db := gormConn()
	var user []User
	db.Raw("SELECT * FROM `users` WHERE NOT subscription IS NULL;").Scan(&user)
	fmt.Println(user[0].Id)
	return user
}
func CheckSub(user []User) {
	fmt.Println(user[0].Id)
	// fmt.Println(user[1].Id)
	for _, u := range user {
		dateDiff := calculateTimeDiff(*u.Subscription)
		if dateDiff <= 3 && dateDiff > 0 {
			// SendAlmostRanOutEmail(u)
		} else if dateDiff <= 0 {
			db := gormConn()
			// SendSubsRanOutEmail(u)
			db.Model(&u).Update("subscription", nil)
		}
	}
}

func calculateTimeDiff(subEndDate string) int {
	timeFormat := "2006-01-02"
	a, err := time.Parse(timeFormat, subEndDate)
	currDate := time.Now().Local()
	if err != nil {
		panic(err)
	}
	sisa := math.Ceil(a.Sub(currDate).Hours() / 24)
	return int(sisa)
}
