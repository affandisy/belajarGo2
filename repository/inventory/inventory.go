package inventory

import (
	"belajarGo2/service/inventory"
	"context"

	"gorm.io/gorm"
)

type (
	GormRepository struct {
		*gorm.DB
	}
)

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		// db.Table("inventories"),
		db.Table("bg_inventories"),
	}
}

func (r *GormRepository) Create(inv inventory.Inventory) (err error) {
	ctx := context.Background()
	return r.DB.WithContext(ctx).Create(&inv).Error
}

func (r *GormRepository) ReadAll(page int, limit int) (invs []inventory.Inventory, err error) {
	ctx := context.Background()
	// r.DB.WithContext(ctx).Offset((page - 1) * limit).Limit(limit).Find(&invs)
	r.DB.WithContext(ctx).Order("code DESC").Offset((page - 1) * limit).Limit(limit).Find(&invs)
	return
}

func (r *GormRepository) ReadByCode(code string) (inv inventory.Inventory, err error) {
	ctx := context.Background()
	r.DB.WithContext(ctx).First(&inv, "code = ?", code)
	return
}

func (r *GormRepository) Update(inv inventory.Inventory) (err error) {
	ctx := context.Background()
	return r.DB.WithContext(ctx).Where("code = ?", inv.Code).Save(inv).Error
}

func (r *GormRepository) Delete(code string) (err error) {
	ctx := context.Background()
	return r.DB.WithContext(ctx).Where("code = ?", code).Delete(inventory.Inventory{}).Error
}
