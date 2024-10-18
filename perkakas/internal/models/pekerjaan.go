package models

import "github.com/tigapilarmandiri/perkakas/common/db"

type Pekerjaan struct {
	db.Model

	Kode   string `json:"kode" gorm:"type:varchar;size:15;not null" validate:"required"`
	Nama   string `json:"nama" gorm:"type:varchar;size:200;not null" validate:"required"`
	Active *bool  `json:"active" validate:"required"`

	ParentID *string    `json:"parent_id" gorm:"index, type:uuid" validate:"omitempty,uuid4"`
	Parent   *Pekerjaan `json:"parent,omitempty" gorm:"foreignKey:ParentID" validate:"-"`
}
