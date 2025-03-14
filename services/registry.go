package services

import (
	"user-service/repositories"
	"user-service/services/user"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() user.IUserService
}

func NewServiceRegistry(repository repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{
		repository: repository,
	}
}
func (r *Registry) GetUser() user.IUserService {
	return user.NewUserService(r.repository)
}
