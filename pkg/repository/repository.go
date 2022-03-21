package repository

import (
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/pvs9/everphone-test-task/pkg/repository/pgx"
	"github.com/uptrace/bun"
)

type Employee interface {
	CreateByNames(names []string) ([]christmas.Employee, error)
	DetachTagsFromManyByIds(ids []int64) error
	GetAll() ([]christmas.Employee, error)
	GetById(id int64) (*christmas.Employee, error)
	GetByNames(names []string, createMissing bool) ([]christmas.Employee, error)
	Update(employee christmas.Employee) (christmas.Employee, error)
	SyncTags(employee christmas.Employee, tags []christmas.Tag) (christmas.Employee, error)
}

type Gift interface {
	CreateByNames(names []string) ([]christmas.Gift, error)
	DetachTagsFromManyByIds(ids []int64) error
	GetAll() ([]christmas.Gift, error)
	GetById(id int64) (*christmas.Gift, error)
	GetByNames(names []string, createMissing bool) ([]christmas.Gift, error)
	Update(gift christmas.Gift) (christmas.Gift, error)
	SyncTags(gift christmas.Gift, tags []christmas.Tag) (christmas.Gift, error)
}

type Tag interface {
	AttachManyEmployees(employeeTags []christmas.EmployeeTag) ([]christmas.EmployeeTag, error)
	AttachManyGifts(giftTags []christmas.GiftTag) ([]christmas.GiftTag, error)
	CreateByNames(names []string) ([]christmas.Tag, error)
	GetByNames(names []string, createMissing bool) ([]christmas.Tag, error)
}

type Repository struct {
	Employee
	Gift
	Tag
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		Employee: pgx.NewEmployeePGX(db),
		Gift:     pgx.NewGiftPGX(db),
		Tag:      pgx.NewTagPGX(db),
	}
}
