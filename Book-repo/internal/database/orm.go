package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Document struct {
	gorm.Model
	TicketID string `gorm:"type:varchar(100);unique_index"`
	Content  string `gorm:"type:varchar(100)"`
	Title    string `gorm:"type:varchar(100)"`
	Author   string `gorm:"type:varchar(100)"`
	Topic    string `gorm:"type:varchar(100)"`
	Bookmark int    `gorm:"type:int"`
}
type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (Document) TableName() string {
	return "BookRepository"
}
func Init(dialect, host, port, dbname, pass string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dialect,
		host, port, dbname, pass))
	if err != nil {
		log.Println("Failed to connect to database")
	} else {
		//db.AutoMigrate(&Document{})
		//	m := Document{TicketID: "1", Content: "book", Title: "Harry Potter", Author: "J.K.Rowling", Topic: "Dobby", Bookmark: 0}
		//db.Create(&m)
		//db.Save(&m)
		//doc := []Document{}
		//db.Find(&doc)
		//log.Println(doc)
	}
	return db, err
	//defer db.Close()
}
