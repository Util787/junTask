package service

import "github.com/Util787/junTask/internal/repository"

type UserService interface {
}

type Service struct {
	UserService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
