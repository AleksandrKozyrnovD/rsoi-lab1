package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"lab1/models"
)

var (
	NotFoundError = errors.New("record not found")
)

type PersonRepository interface {
	GetAll() ([]models.Person, error)
	GetById(id int32) (*models.Person, error)
	AddPerson(person *models.Person) (int32, error)
	DeletePerson(id int32) error
	UpdatePerson(id int32, req *models.PersonRequest) (*models.Person, error)
}

type personRepository struct {
	DB *gorm.DB
}

func NewPersonRepository(db *gorm.DB) PersonRepository {
	return &personRepository{DB: db}
}

func (r *personRepository) GetAll() ([]models.Person, error) {
	persons := make([]models.Person, 0)
	err := r.DB.Find(&persons).Error
	if err != nil {
		return nil, err
	}

	return persons, nil
}

func (r *personRepository) GetById(id int32) (*models.Person, error) {
	var person models.Person
	if err := r.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("person %d: %w", id, NotFoundError)
		}
		return nil, fmt.Errorf("person %d: %w", id, err)
	}

	return &person, nil
}

func (r *personRepository) AddPerson(person *models.Person) (int32, error) {
	if person == nil {
		return 0, errors.New("person is nil")
	}

	err := r.DB.Create(person).Error

	if err != nil {
		return 0, err
	}

	return person.Id, nil
}

func (r *personRepository) DeletePerson(id int32) error {
	result := r.DB.Delete(&models.Person{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return NotFoundError
	}

	return nil
}

func (r *personRepository) UpdatePerson(id int32, req *models.PersonRequest) (*models.Person, error) {
	if req == nil {
		return nil, errors.New("req is nil")
	}

	person, err := r.GetById(id)
	if err != nil {
		return nil, NotFoundError
	}

	err = r.DB.Model(person).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":    req.Name,
			"address": req.Address,
			"age":     req.Age,
			"work":    req.Work,
		}).Error

	if err != nil {
		return nil, err
	}
	person.Name = req.Name
	person.Address = req.Address
	person.Age = req.Age
	person.Work = req.Work

	return person, nil
}
