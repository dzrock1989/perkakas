package pekerjaan

import (
	"context"

	"github.com/tigapilarmandiri/perkakas/internal/models"
	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewPekerjaanRepository(DB *gorm.DB) *repository {
	return &repository{
		DB: DB,
	}
}

func (repo *repository) Create(ctx context.Context, direktorat *models.Pekerjaan) (err error) {
	if err = repo.DB.Create(direktorat).Error; err != nil {
		return
	}

	return
}

func (repo *repository) Update(ctx context.Context, direktorat models.Pekerjaan) (err error) {
	if rowAffected := repo.DB.
		Where("id = ?", direktorat.ID).
		Updates(direktorat).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}

func (repo *repository) Delete(ctx context.Context, Uuid string) (err error) {
	if rowAffected := repo.DB.
		Where("id = ?", Uuid).
		Delete(&models.Pekerjaan{}).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}
