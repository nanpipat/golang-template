package repo

import (
	"github.com/nanpipat/golang-template-hexagonal/helper"
	"github.com/nanpipat/golang-template-hexagonal/internal/core/domain"
	"github.com/nanpipat/golang-template-hexagonal/utils"

	"gorm.io/gorm"
)

type IUserRepositoryInterface interface {
	Create(payload domain.User) (*domain.User, error)
	Pagination(pageOptions *helper.PageOptions) (*Pagination[domain.User], error)
	Get(id string) (*domain.User, error)
	Update(id string, payload *domain.User) (*domain.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepositoryInterface {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(payload domain.User) (*domain.User, error) {
	err := New[domain.User](r.db).Create(&payload)
	if err != nil {
		return nil, err
	}
	return r.Get(payload.ID)
}

func (r *userRepository) Pagination(pageOptions *helper.PageOptions) (*Pagination[domain.User], error) {
	db := New[domain.User](r.db)

	result, err := db.
		Pagination(pageOptions)
	if err != nil {
		return nil, err
	}

	// return PaginationMap(result, func(p domain.User) domain.User {
	// 	return *views.NewCMSManual(&p)
	// }), nil

	return result, nil
}

func (r *userRepository) Get(id string) (*domain.User, error) {
	return New[domain.User](r.db).FindOne("id = ?", id)
}

func (r *userRepository) Update(id string, payload *domain.User) (*domain.User, error) {
	user, err := r.Get(id)
	if err != nil {
		return nil, err
	}

	err = utils.Copy(&user, &payload)
	if err != nil {
		return nil, err
	}

	err = New[domain.User](r.db).Update(&payload)
	if err != nil {
		return nil, err
	}

	return r.Get(id)
}

func (r *userRepository) Delete(id string) error {
	err := New[domain.User](r.db).Delete("id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
