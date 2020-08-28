package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zdarovich/win-lose-api/models"
	"log"
)

var (
	db *gorm.DB
)
func InitDB() {
	var err error
	db, err = gorm.Open("postgres", "host=postgres port=5432 user=user dbname=winlose password=pass sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Transaction{}, &models.User{})
	var user models.User
	db.FirstOrCreate(&user, models.User{Name: "TestUser", Balance:0})
	db.Model(user).Update("balance", 0)
}

func GetDB() *gorm.DB {
	return db
}