package main

import (
	"os"
	"ping-util/internal/mx"
	"ping-util/internal/record"
	"ping-util/pkg/gormclient"
)

func main() {
	gorm := gormclient.NewClient(os.Getenv("POSTGRES_DSN"))
	if err := gorm.AutoMigrate(&mx.DomainRepo{}); err != nil {
		panic(err)
	}

	domainRepo := mx.NewDomainRepo(gorm)
	domainUsecase := mx.NewDomainUsecase(domainRepo)

	recordRepo := record.NewRecordRepo(os.Getenv("JSON_FILE_PATH"))
	_ = record.NewRecordUsecase(recordRepo)

	if err := domainUsecase.EnrichDomainList(); err != nil {
		panic(err)
	}
}
