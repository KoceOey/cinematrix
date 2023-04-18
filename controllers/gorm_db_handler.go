package controllers

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func gormConn() *gorm.DB {
	dsn := "root:@tcp(localhost:3306)/db_tubes_pbp"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
