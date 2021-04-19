package database

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"Book-repo/internal"
	"Book-repo/internal/database"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"

	"github.com/lithammer/shortuuid/v3"
)

type dbService struct{}

func NewService() Service { return &dbService{} }

func (d *dbService) Add(_ context.Context, doc *internal.Document) (string, error) {
	newTicketID := shortuuid.New()
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", database.DefaultHost, database.DefaultPort, database.DefaultDBUser, database.DefaultDatabase,
		database.DefaultPassword))

	defer func() {
		err := db.Close()
		if err != nil {
			logger.Log("ERROR::Failed to close the database connection ", err.Error())
		}
	}()
	if err != nil {
		logger.Log(fmt.Sprintf("FATAL: failed to load db with error: %s", err.Error()))
	}

	m := database.Document{TicketID: newTicketID, Content: doc.Content, Title: doc.Title, Author: doc.Author, Topic: doc.Topic, Bookmark: 0}
	db.Create(&m)
	return newTicketID, nil
}

func (d *dbService) Get(_ context.Context, filters ...internal.Filter) ([]internal.Document, error) {

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", database.DefaultHost, database.DefaultPort, database.DefaultDBUser, database.DefaultDatabase,
		database.DefaultPassword))

	defer func() {
		err := db.Close()
		if err != nil {
			logger.Log("ERROR::Failed to close the database connection ", err.Error())
		}
	}()
	if err != nil {
		logger.Log(fmt.Sprintf("FATAL: failed to load db with error: %s", err.Error()))
	}
	//logger.Log(filters[0].Key, filters[0].Value, len(filters))
	var result []database.Document
	var querycond database.Document
	result = make([]database.Document, 5, 10)
	for _, f := range filters {
		//db.Where(f.Key, "=?", f.Value).Find(&result[i])

		if f.Key == "Title" {
			querycond.Title = f.Value
			//db.Where(&database.Document{Title: f.Value}).First(&result)
		} else if f.Key == "Author" {
			querycond.Author = f.Value
			//	db.Where(&database.Document{Author: f.Value}).First(&result)
		} else if f.Key == "Content" {
			querycond.Content = f.Value
			//	db.Where(&database.Document{Content: f.Value}).First(&result)
		} else if f.Key == "Topic" {
			querycond.Topic = f.Value
			//	db.Where(&database.Document{Topic: f.Value}).First(&result)
		}
		//db.Where("? = ?", Key, Val).First(&result[i])
		//db.Find(&result[i], "?=?", Key, Val)
		//	db.Exec("SELECT * FROM documents WHERE Title = ?", f.Value).Scan(&result[i])
		//db.Find(&result[i], "Title=?", "Harry Potter")
	}
	db.Where(&querycond).Find(&result)
	/*doc1 := internal.Document{
		Content: "book",
		Title:   "Harry Potter",
		Author:  "J.K. Rowling",
		Topic:   "Fiction and Magic",
	} */
	var docs []internal.Document
	docs = make([]internal.Document, len(result), len(result))
	for j, res := range result {
		docs[j] = internal.Document{
			Content:  res.Content,
			Title:    res.Title,
			Author:   res.Author,
			Topic:    res.Topic,
			Bookmark: res.Bookmark,
		}
	}
	return docs, nil
}

func (d *dbService) Update(_ context.Context, ticketId string, doc *internal.Document) (int, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", database.DefaultHost, database.DefaultPort, database.DefaultDBUser, database.DefaultDatabase,
		database.DefaultPassword))

	defer func() {
		err := db.Close()
		if err != nil {
			logger.Log("ERROR::Failed to close the database connection ", err.Error())
		}
	}()
	if err != nil {
		logger.Log(fmt.Sprintf("FATAL: failed to load db with error: %s", err.Error()))
	}
	var result database.Document
	db.Model(&result).Where(&database.Document{TicketID: ticketId}).Update("Bookmark", doc.Bookmark)
	return http.StatusOK, nil
}

func (d *dbService) Remove(_ context.Context, ticketId string) (int, error) {
	return http.StatusOK, nil
}

//func (d *dbService) Validate(_ context.Context, doc *internal.Document) (bool, error) {
//	return true, nil
//}

func (d *dbService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
