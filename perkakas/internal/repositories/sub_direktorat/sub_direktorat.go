package subdirektorat

import (
	"context"

	"github.com/tigapilarmandiri/perkakas/internal/models"
	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewSubDirektoratRepository(DB *gorm.DB) *repository {
	return &repository{
		DB: DB,
	}
}

func (repo *repository) Create(ctx context.Context, subDirektorat *models.SubDirektorat) (err error) {
	if err = repo.DB.Create(subDirektorat).Error; err != nil {
		return
	}

	return
}

func (repo *repository) Update(ctx context.Context, subDirektorat models.SubDirektorat) (err error) {
	if rowAffected := repo.DB.
		Model(&models.SubDirektorat{}).Where("id = ?", subDirektorat.ID).
		Updates(subDirektorat).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}

func (repo *repository) Delete(ctx context.Context, Uuid string) (err error) {
	if rowAffected := repo.DB.Where("id = ?", Uuid).Delete(&models.SubDirektorat{}).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}
