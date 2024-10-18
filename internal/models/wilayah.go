package models

import (
	"database/sql/driver"

	"github.com/tigapilarmandiri/perkakas/common/db"
)

type JenisWilayah string

const (
	NASIONAL  JenisWilayah = "NASIONAL"
	PROVINSI  JenisWilayah = "PROVINSI"
	KABKOTA   JenisWilayah = "KABKOTA"
	KECAMATAN JenisWilayah = "KECAMATAN"
	DESALURAH JenisWilayah = "DESALURAH"
)

func (ct *JenisWilayah) Scan(value interface{}) error {
	*ct = JenisWilayah(value.(string))
	return nil
}

func (ct JenisWilayah) Value() (driver.Value, error) {
	return string(ct), nil
}

type Wilayah struct {
	db.Model

	ParentID *string  `json:"parent_id" gorm:"index, type:uuid"`
	Parent   *Wilayah `json:"parent,omitempty" gorm:"foreignKey:ParentID"`

	Kode      string       `json:"kode" gorm:"type:varchar;size:15" validate:"required"`
	Nama      string       `json:"nama" gorm:"type:varchar;size:255;not null" validate:"required"`
	Koordinat string       `json:"koordinat" gorm:"type:varchar;size:100"`
	Jenis     JenisWilayah `json:"jenis" gorm:"type:jenis_wilayah;not null;default:DESALURAH" validate:"required" filter:"jenis"`
	Active    *bool        `json:"active" validate:"required"`

	Kepolisians []*Kepolisian `json:"kepolisians,omitempty" gorm:"many2many:kepolisian_has_wilayah"`
}
