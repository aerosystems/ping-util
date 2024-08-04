package mx

import (
	"gorm.io/gorm"
	"net"
	"time"
)

const defaultBulkSize = 1000

type Repo struct {
	db *gorm.DB
}

func NewDomainRepo(db *gorm.DB) *Repo {
	return &Repo{
		db,
	}
}

type DomainRepo struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"unique;index"`
	Type      string    `gorm:"index"`
	Ip        string    `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (d *DomainRepo) toModel() *DomainModel {
	return &DomainModel{
		Name: d.Name,
		Type: DomainType{d.Type},
		Ip:   net.ParseIP(d.Ip),
	}
}

func modelToDomainRepo(d *DomainModel) *DomainRepo {
	return &DomainRepo{
		Name: d.Name,
		Type: d.Type.String(),
		Ip:   d.Ip.String(),
	}
}

func (r *Repo) Create(d *DomainModel) error {
	return r.db.Create(modelToDomainRepo(d)).Error
}

func (r *Repo) FindByName(name string) (*DomainModel, error) {
	var d DomainRepo
	err := r.db.Where("name = ?", name).First(&d).Error
	if err != nil {
		return nil, err
	}
	return d.toModel(), nil
}

func (r *Repo) CreateDomainList(domains []*DomainModel) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var bulkDomains []*DomainRepo
	for i, d := range domains {
		if len(bulkDomains) < defaultBulkSize {
			bulkDomains = append(bulkDomains, modelToDomainRepo(d))
			if i != len(domains)-1 {
				continue
			}
		}
		if err := tx.Create(&bulkDomains).Error; err != nil {
			tx.Rollback()
		}
		bulkDomains = nil
	}
	return tx.Commit().Error
}

func (r *Repo) Update(d *DomainModel) error {
	domain := modelToDomainRepo(d)
	return r.db.Model(&DomainRepo{}).Where("name = ?", d.Name).Updates(domain).Error
}

func (r *Repo) FindWithoutAddress() ([]*DomainModel, error) {
	var domainRepos []*DomainRepo
	if err := r.db.Find(&domainRepos).Where("ip = ?", "").Error; err != nil {
		return nil, err
	}
	var domainModels []*DomainModel
	for _, d := range domainRepos {
		domainModels = append(domainModels, d.toModel())
	}
	return domainModels, nil
}
