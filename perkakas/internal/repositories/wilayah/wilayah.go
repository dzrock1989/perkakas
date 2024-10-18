package wilayah

import (
	"context"

	"github.com/tigapilarmandiri/perkakas/internal/models"
	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewWilayahRepository(DB *gorm.DB) *repository {
	return &repository{
		DB: DB,
	}
}

func (repo *repository) Create(ctx context.Context, wilayah *models.Wilayah) (err error) {
	if err = repo.DB.Create(wilayah).Error; err != nil {
		return
	}

	return
}

func (repo *repository) Update(ctx context.Context, wilayah models.Wilayah) (err error) {
	if rowAffected := repo.DB.
		Model(&models.Wilayah{}).Where("id = ?", wilayah.ID).
		Updates(wilayah).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}

func (repo *repository) Delete(ctx context.Context, Uuid string) (err error) {
	if rowAffected := repo.DB.Where("id = ?", Uuid).Delete(&models.Wilayah{}).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}
