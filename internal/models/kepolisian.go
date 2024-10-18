package models

import (
	"database/sql/driver"

	"github.com/tigapilarmandiri/perkakas/common/db"
)

type JenisKepolisian string

const (
	MABES  JenisKepolisian = "MABES"
	POLDA  JenisKepolisian = "POLDA"
	POLRES JenisKepolisian = "POLRES"
	POLSEK JenisKepolisian = "POLSEK"
)

func (ct *JenisKepolisian) Scan(value interface{}) error {
	*ct = JenisKepolisian(value.(string))
	return nil
}

func (ct JenisKepolisian) Value() (driver.Value, error) {
	return string(ct), nil
}

type Kepolisian struct {
	db.Model

	Kode      string          `json:"kode" gorm:"type:varchar;size:15;not null" validate:"required"`
	Nama      string          `json:"nama" gorm:"type:varchar;size:255;not null" validate:"required"`
	Koordinat string          `json:"koordinat" gorm:"type:varchar;size:100;not null" validate:"required"`
	Jenis     JenisKepolisian `json:"jenis" gorm:"type:jenis_kepolisian;not null;default:POLSEK" validate:"required" filter:"jenis"`
	Active    *bool           `json:"active" validate:"required"`

	ParentID *string     `json:"parent_id" gorm:"type:uuid"`
	Parent   *Kepolisian `json:"parent,omitempty" gorm:"foreignKey:ParentID"`

	Wilayahs []*Wilayah `json:"wilayah,omitempty" gorm:"many2many:kepolisian_has_wilayah"`
}
