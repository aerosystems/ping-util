package mx

import (
	"context"
	"golang.org/x/sync/semaphore"
	"log"
	"net"
	"time"
)

type Repository interface {
	Create(d *DomainModel) error
	FindByName(name string) (*DomainModel, error)
	CreateDomainList(domains []*DomainModel) error
	Update(d *DomainModel) error
	FindWithoutAddress() ([]*DomainModel, error)
}

type Usecase struct {
	repo Repository
}

func NewDomainUsecase(repo Repository) *Usecase {
	return &Usecase{
		repo,
	}
}

func (u *Usecase) AddDomainList(domains []string) error {
	var domainModels []*DomainModel
	for _, d := range domains {
		domainModels = append(domainModels, &DomainModel{
			Name: d,
			Type: UnknownDomainType,
			Ip:   nil,
		})
	}
	return u.repo.CreateDomainList(domainModels)
}

func (u *Usecase) EnrichDomainList() error {
	domains, err := u.repo.FindWithoutAddress()
	if err != nil {
		return err
	}
	sm := semaphore.NewWeighted(1000)
	for _, d := range domains {
		if err := sm.Acquire(context.Background(), 1); err != nil {
			return err
		}
		go func(domain *DomainModel) {
			defer sm.Release(1)
			u.enrichDomain(domain)
		}(d)
	}
	return nil
}

func (u *Usecase) enrichDomain(domain *DomainModel) {
	ip, err := lookup(domain.Name)
	if err != nil {
		log.Printf("failed to lookup %s: %v", domain.Name, err)
	}
	domain.Ip = ip
	if err := u.repo.Update(domain); err != nil {
		log.Printf("failed to update %s: %v", domain.Name, err)
	}
}

func lookup(domain string) (net.IP, error) {
	log.Println("looking up: ", domain)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, domain)
	if err != nil {
		return nil, err
	}
	ip := ips[0].IP
	return ip, nil
}
