package service

import (
	"errors"
	"lab1/models"
	"lab1/repository"
)

var (
	ErrNotFound = errors.New("person not found")
)

type PersonService interface {
	List() ([]models.Person, error)
	GetByID(id int32) (models.Person, error)
	Create(p models.Person) (models.Person, error)
	Update(id int32, p models.PersonRequest) (models.Person, error)
	Delete(id int32) error
}

type personService struct {
	personRepository repository.PersonRepository
}

func NewPersonService(personRepository repository.PersonRepository) PersonService {
	return &personService{
		personRepository: personRepository,
	}
}

func (s *personService) List() ([]models.Person, error) {
	persons, err := s.personRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return persons, nil
}

func (s *personService) GetByID(id int32) (models.Person, error) {
	person, err := s.personRepository.GetById(id)
	if err != nil {
		if errors.Is(err, repository.NotFoundError) {
			return models.Person{}, ErrNotFound
		}
		return models.Person{}, err
	}
	return *person, nil
}

func (s *personService) Create(p models.Person) (models.Person, error) {
	id, err := s.personRepository.AddPerson(&p)
	if err != nil {
		return models.Person{}, err
	}
	p.Id = id
	return p, nil
}

func (s *personService) Update(id int32, p models.PersonRequest) (models.Person, error) {
	person, err := s.personRepository.UpdatePerson(id, &p)
	if err != nil {
		if errors.Is(err, repository.NotFoundError) {
			return models.Person{}, ErrNotFound
		}
		return models.Person{}, err
	}
	return *person, nil
}

func (s *personService) Delete(id int32) error {
	err := s.personRepository.DeletePerson(id)
	if err != nil {
		if errors.Is(err, repository.NotFoundError) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
