package service

import (
	"errors"
	"testing"

	"lab1/models"
	"lab1/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPersonRepository — мок-реализация интерфейса repository.PersonRepository
type MockPersonRepository struct {
	mock.Mock
}

func (m *MockPersonRepository) GetAll() ([]models.Person, error) {
	args := m.Called()
	return args.Get(0).([]models.Person), args.Error(1)
}

func (m *MockPersonRepository) GetById(id int32) (*models.Person, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Person), args.Error(1)
}

func (m *MockPersonRepository) AddPerson(person *models.Person) (int32, error) {
	args := m.Called(person)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockPersonRepository) DeletePerson(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPersonRepository) UpdatePerson(id int32, req *models.PersonRequest) (*models.Person, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Person), args.Error(1)
}

// Проверка, что мок реализует интерфейс
var _ repository.PersonRepository = (*MockPersonRepository)(nil)

// ---------- List ----------

func TestPersonService_List_Success(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	expected := []models.Person{
		{Id: 1, Name: "Alice", Age: 30, Address: "123 Main St", Work: "Engineer"},
		{Id: 2, Name: "Bob", Age: 25, Address: "456 Elm St", Work: "Designer"},
	}
	mockRepo.On("GetAll").Return(expected, nil)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.List()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_List_Error(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	dbErr := errors.New("db error")
	mockRepo.On("GetAll").Return([]models.Person{}, dbErr)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.List()

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

// ---------- GetByID ----------

func TestPersonService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	person := &models.Person{Id: 1, Name: "Alice", Age: 30, Address: "123 Main St", Work: "Engineer"}
	mockRepo.On("GetById", int32(1)).Return(person, nil)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, *person, result)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	mockRepo.On("GetById", int32(1)).Return(nil, repository.NotFoundError)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.GetByID(1)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.Equal(t, models.Person{}, result)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_GetByID_Error(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	dbErr := errors.New("db error")
	mockRepo.On("GetById", int32(1)).Return(nil, dbErr)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.GetByID(1)

	assert.ErrorIs(t, err, dbErr)
	assert.Equal(t, models.Person{}, result)
	mockRepo.AssertExpectations(t)
}

// ---------- Create ----------

func TestPersonService_Create_Success(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	input := models.Person{Name: "Alice", Age: 30, Address: "123 Main St", Work: "Engineer"}
	expectedID := int32(1)
	mockRepo.On("AddPerson", mock.AnythingOfType("*models.Person")).Return(expectedID, nil)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.Create(input)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.Id)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Age, result.Age)
	assert.Equal(t, input.Address, result.Address)
	assert.Equal(t, input.Work, result.Work)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_Create_Error(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	input := models.Person{Name: "Alice", Age: 30}
	dbErr := errors.New("db error")
	mockRepo.On("AddPerson", mock.AnythingOfType("*models.Person")).Return(int32(0), dbErr)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.Create(input)

	assert.ErrorIs(t, err, dbErr)
	assert.Equal(t, models.Person{}, result)
	mockRepo.AssertExpectations(t)
}

// ---------- Update ----------

func TestPersonService_Update_Success(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	req := models.PersonRequest{Name: "Alice Updated", Age: 31, Address: "New Address", Work: "Senior Engineer"}
	updatedPerson := &models.Person{Id: 1, Name: "Alice Updated", Age: 31, Address: "New Address", Work: "Senior Engineer"}
	mockRepo.On("UpdatePerson", int32(1), mock.AnythingOfType("*models.PersonRequest")).Return(updatedPerson, nil)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.Update(1, req)

	assert.NoError(t, err)
	assert.Equal(t, *updatedPerson, result)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_Update_Error(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	req := models.PersonRequest{Name: "Alice Updated", Age: 31}
	dbErr := errors.New("db error")
	mockRepo.On("UpdatePerson", int32(1), mock.AnythingOfType("*models.PersonRequest")).Return(nil, dbErr)

	svc := &personService{personRepository: mockRepo}
	result, err := svc.Update(1, req)

	assert.ErrorIs(t, err, dbErr)
	assert.Equal(t, models.Person{}, result)
	mockRepo.AssertExpectations(t)
}

// ---------- Delete ----------

func TestPersonService_Delete_Success(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	mockRepo.On("DeletePerson", int32(1)).Return(nil)

	svc := &personService{personRepository: mockRepo}
	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_Delete_Error(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	dbErr := errors.New("db error")
	mockRepo.On("DeletePerson", int32(1)).Return(dbErr)

	svc := &personService{personRepository: mockRepo}
	err := svc.Delete(1)

	assert.ErrorIs(t, err, dbErr)
	mockRepo.AssertExpectations(t)
}
