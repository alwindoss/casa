package repository

import "gorm.io/gorm"

type UserRepositoy interface {
}

func NewGORMUserRepository(db *gorm.DB) UserRepositoy {
	return &gormUserRepository{
		db: db,
	}
}

type gormUserRepository struct {
	db *gorm.DB
}
