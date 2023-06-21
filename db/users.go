package db

import (
	"log"

	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Token       string `gorm:"unique"`
	Email       string
	WorkspaceId int
}

type DataBase struct {
	orm *gorm.DB
}

func NewDatabase() *DataBase {
	orm, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	orm.AutoMigrate(&User{})

	return &DataBase{orm: orm}
}

func (db *DataBase) CreateUser(token string, workspaceId int) {
	log.Printf("[DEBUG] db.CreateUser")
	user := &User{Token: token, WorkspaceId: workspaceId}
	db.orm.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
}
