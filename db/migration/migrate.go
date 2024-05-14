package main

import (
	"fmt"
	"log"
	"os"

	"simple-users-tasks-service/db/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.User{}, &model.Task{})

	db.Create(&model.User{Email: "test@gmail.com", Name: "Alice"})
	db.Create(&model.Task{Title: "todo1", UserID: 1})

}
