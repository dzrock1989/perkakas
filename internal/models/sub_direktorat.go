package models

import "github.com/tigapilarmandiri/perkakas/common/db"

type SubDirektorat struct {
	db.Model

	Nama         string      `json:"nama" gorm:"type:varchar;size:200;not null" validate:"required"`
	DirektoratId *string     `json:"direktorat_id" gorm:"type:uuid"`
	Direktorat   *Direktorat `json:"direktorat,omitempty"`
}
