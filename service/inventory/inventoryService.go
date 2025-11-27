package inventory

type service struct {
	repo Repository
}

type Service interface {
	Create(inv Inventory) (err error)
	GetAll(page int, limit int) (invs []Inventory, err error)
	GetByCode(code string) (inv Inventory, err error)
	Update(code string) (err error)
	Delete(code string) (err error)
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) Create(inv Inventory) (err error) {
	return s.repo.Create(inv)
}

func (s *service) GetAll(page int, limit int) (invs []Inventory, err error) {
	return s.repo.ReadAll(page, limit)
}

func (s *service) GetByCode(code string) (inv Inventory, err error) {
	return s.repo.ReadByCode(code)
}

func (s *service) Update(code string) (err error) {
	return s.repo.Update(code)
}

func (s *service) Delete(code string) (err error) {
	return s.repo.Delete(code)
}
