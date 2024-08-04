package gormclient

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

const countAttempts = 5

var (
	instance *gorm.DB
	once     sync.Once
)

func NewClient(postgresDSN string) *gorm.DB {
	once.Do(func() {
		count := 0
		for {
			db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
			})
			if err != nil {
				log.Println("PostgreSQL not ready...")
				count++
			} else {
				log.Println("Connected to database!")
				instance = db
				return
			}
			if count > countAttempts {
				log.Printf("Failed to connect to database after %d attempts", countAttempts)
				return
			}
			// Exponential backoff
			exponential := 2 << count
			log.Printf("Retrying in %d seconds...", exponential)
			<-time.After(time.Duration(exponential) * time.Second)
		}
	})
	return instance
}
