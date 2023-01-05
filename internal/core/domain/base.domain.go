package domain

import (
	"time"

	"github.com/nanpipat/golang-template-hexagonal/utils"
)

type BaseModel struct {
	ID        string     `json:"id" gorm:"column:id;primary_key"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at"`
}

func NewBaseModel() BaseModel {
	return BaseModel{
		ID:        utils.GetUUID(),
		CreatedAt: utils.GetCurrentDateTime(),
		UpdatedAt: utils.GetCurrentDateTime(),
	}
}
