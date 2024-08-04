package record

type Repository interface {
	GetAll() ([]Domain, error)
}

type Usecase struct {
	repo Repository
}

func NewRecordUsecase(repo Repository) *Usecase {
	return &Usecase{
		repo,
	}
}

func (u *Usecase) GetDomains() ([]string, error) {
	records, err := u.repo.GetAll()
	if err != nil {
		return []string{}, err
	}
	var domains []string
	for _, r := range records {
		domains = append(domains, r.Name)
	}
	return domains, nil
}
