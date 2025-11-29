package inventory

type Repository interface {
	Create(inv Inventory) (err error)
	ReadAll(page int, limit int) (invs []Inventory, err error)
	ReadByCode(code string) (inv Inventory, err error)
	// Update(code string) (err error)
	Update(inv Inventory) (err error)
	Delete(code string) (err error)
}
