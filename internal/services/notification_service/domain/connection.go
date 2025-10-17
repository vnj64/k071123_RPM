package domain

import (
	"gorm.io/gorm"
)

type Connection interface {
	DB() *gorm.DB
}
