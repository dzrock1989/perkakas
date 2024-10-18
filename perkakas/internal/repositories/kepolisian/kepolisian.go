package kepolisian

import (
	"context"
	"fmt"

	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewKepolisianRepository(DB *gorm.DB) *Repository {
	return &Repository{
		DB: DB,
	}
}

func (repo *Repository) Tx(f func(*Repository) error) (err error) {
	tx := NewKepolisianRepository(repo.DB.Begin())

	err = f(tx)

	defer func() {
		if p := recover(); p != nil {
			util.Log.Error().Msg(fmt.Sprint(p))
			tx.DB.Rollback()
			return
		}

		if err != nil {
			tx.DB.Rollback()
			return
		}

		tx.DB.Commit()
	}()

	return
}

func (repo *Repository) Create(ctx context.Context, kepolisian *models.Kepolisian) (err error) {
	if err = repo.DB.Create(kepolisian).Error; err != nil {
		return
	}

	return
}

func (repo *Repository) Update(ctx context.Context, kepolisian models.Kepolisian) (err error) {
	if rowAffected := repo.DB.
		Where("id = ?", kepolisian.ID).
		Updates(kepolisian).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}

func (repo *Repository) Delete(ctx context.Context, Uuid string) (err error) {
	if rowAffected := repo.DB.
		Where("id = ?", Uuid).
		Delete(&models.Kepolisian{}).RowsAffected; rowAffected == 0 {
		err = gorm.ErrRecordNotFound

		return
	}

	return
}

func (repo *Repository) GetKepolisianWilayah(ctx context.Context, uuid string, results *[]models.Wilayah) (err error) {
	err = repo.DB.Raw("select from wilayahs w left join kepolisian_has_wilayah khw on khw.wilayah_id = w.id where khw.kepolisian_id = ?", uuid).Scan(results).Error
	return
}

func (repo *Repository) DeleteKepolisianIdOnTableJoin(ctx context.Context, uuid string) (err error) {
	err = repo.DB.Exec("delete from kepolisian_has_wilayah where kepolisian_id = ?", uuid).Error
	return
}

func (repo *Repository) SyncKepolisianHasWilayah(ctx context.Context, uuid string, datas []string) (err error) {
	q := "insert into kepolisian_has_wilayah (kepolisian_id, wilayah_id) values "
	q += fmt.Sprintf("('%s', '%s')", uuid, datas[0])
	for _, v := range datas[1:] {
		q += fmt.Sprintf(", ('%s', '%s')", uuid, v)
	}

	err = repo.DB.Exec(q).Error
	return
}
