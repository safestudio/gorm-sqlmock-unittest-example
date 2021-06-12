package repository

import (
	"github.com/Rosaniline/gorm-ut/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Get(id uuid.UUID) (*model.Person, error)
	Create(id uuid.UUID, name string) error
}

type repo struct {
	DB *gorm.DB
}

func (p *repo) Create(id uuid.UUID, name string) error {
	person := &model.Person{
		ID:   id,
		Name: name,
	}

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(person).Error; err != nil {
			// return any error will rollback
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return err
}

func (p *repo) Get(id uuid.UUID) (*model.Person, error) {
	person := new(model.Person)

	err := p.DB.Where("id = ?", id).Find(person).Error

	return person, err
}

func CreateRepository(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}
