package services

import (
	"github.com/nanpipat/golang-template-hexagonal/helper"
	"github.com/nanpipat/golang-template-hexagonal/internal/core/domain"
	"github.com/nanpipat/golang-template-hexagonal/internal/repo"
)

type IUserServiceInterface interface {
	Create(payload domain.User) (*domain.User, error)
	Pagination(pageOptions *helper.PageOptions) (*repo.Pagination[domain.User], error)
	Get(id string) (*domain.User, error)
	Update(id string, payload *domain.User) (*domain.User, error)
	Delete(id string) error
}

type userService struct {
	repo repo.IUserRepositoryInterface
}

func NewUserService(repo repo.IUserRepositoryInterface) IUserServiceInterface {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(payload domain.User) (*domain.User, error) {
	return s.repo.Create(payload)
}

func (s *userService) Pagination(pageOptions *helper.PageOptions) (*repo.Pagination[domain.User], error) {
	return s.repo.Pagination(pageOptions)
}

func (s *userService) Get(id string) (*domain.User, error) {
	return s.repo.Get(id)
}

func (s *userService) Update(id string, payload *domain.User) (*domain.User, error) {
	return s.repo.Update(id, payload)
}

func (s *userService) Delete(id string) error {
	return s.repo.Delete(id)
}
