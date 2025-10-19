package core

import "gorm.io/gorm"

type Application struct {
	DB *gorm.DB
}
