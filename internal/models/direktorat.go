package models

import "github.com/tigapilarmandiri/perkakas/common/db"

type Direktorat struct {
	db.Model

	Nama string `json:"nama" gorm:"type:varchar;size:200;not null" validate:"required"`
}
