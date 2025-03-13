package repositories

import (
	"user-service/repositories/user"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() user.IUserRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Registry{db: db}
}

func (r *Registry) GetUser() user.IUserRepository {
	return user.NewUserRepository(r.db)
}
