package record

import (
	"encoding/json"
	"io"
	"os"
)

type Record struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

func (r Record) toModel() *Domain {
	return &Domain{
		Name: r.Domain,
	}
}

func recordsToModels(records []Record) []Domain {
	var domains []Domain
	for _, r := range records {
		domains = append(domains, *r.toModel())
	}
	return domains
}

type RecordRepo struct {
	filePath string
}

func NewRecordRepo(filePath string) *RecordRepo {
	return &RecordRepo{
		filePath,
	}
}

func (r RecordRepo) GetAll() ([]Domain, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return []Domain{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return []Domain{}, err
	}
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return []Domain{}, err
	}
	return recordsToModels(records), nil
}
