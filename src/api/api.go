package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lab1/models"
	"lab1/service"
)

type V1 struct {
	personService service.PersonService
}

func NewV1(personService service.PersonService) *V1 {
	return &V1{personService: personService}
}

// GET /api/v1/persons
func (a *V1) HandlePersonsGet(c *gin.Context) {
	persons, err := a.personService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	resp := make([]models.PersonResponse, 0, len(persons))
	for _, p := range persons {
		resp = append(resp, toPersonResponse(p))
	}
	c.JSON(http.StatusOK, resp)
}

// POST /api/v1/persons
func (a *V1) HandlePersonsPost(c *gin.Context) {
	var req models.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Message: "Invalid data",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if errs := validatePersonRequest(req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Message: "Invalid data",
			Errors:  errs,
		})
		return
	}

	created, err := a.personService.Create(toServicePerson(req))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Message: err.Error(),
			Errors:  map[string]string{},
		})
		return
	}

	c.Header("Location", "/api/v1/persons/"+strconv.FormatInt(int64(created.Id), 10))
	c.JSON(http.StatusCreated, toPersonResponse(created))
}

// GET /api/v1/persons/{id}
func (a *V1) HandlePersonsGetId(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	person, err := a.personService.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Not found Person for ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, toPersonResponse(person))
}

// DELETE /api/v1/persons/{id}
func (a *V1) HandlePersonsDeleteId(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := a.personService.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Not found Person for ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /api/v1/persons/{id}
func (a *V1) HandlePersonsPatchId(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req models.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Message: "Invalid data",
			Errors:  map[string]string{"body": err.Error()},
		})
		return
	}

	if errs := validatePersonRequest(req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Message: "Invalid data",
			Errors:  errs,
		})
		return
	}

	existing, err := a.personService.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Not found Person for ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	// PATCH: не затираем поля, которые не пришли в запросе
	if req.Name == "" {
		req.Name = existing.Name
	}
	if req.Age == 0 {
		req.Age = existing.Age
	}
	if req.Address == "" {
		req.Address = existing.Address
	}
	if req.Work == "" {
		req.Work = existing.Work
	}

	updated, err := a.personService.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "Not found Person for ID"})
			return
		}
		c.JSON(http.StatusBadRequest, models.ValidationErrorResponse{
			Message: err.Error(),
			Errors:  map[string]string{},
		})
		return
	}

	c.JSON(http.StatusOK, toPersonResponse(updated))
}

// ---------- helpers ----------

func parseID(c *gin.Context) (int32, bool) {
	raw := c.Param("id")
	v, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "invalid id"})
		return 0, false
	}
	return int32(v), true
}

// validatePersonRequest — минимальная валидация по схеме required: [name]
func validatePersonRequest(req models.PersonRequest) map[string]string {
	errs := map[string]string{}
	if req.Name == "" {
		errs["name"] = "name is required"
	}
	return errs
}

// toServicePerson — API-DTO → доменная модель
func toServicePerson(req models.PersonRequest) models.Person {
	return models.Person{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Work:    req.Work,
	}
}

// toPersonResponse — доменная модель → API-DTO
func toPersonResponse(p models.Person) models.PersonResponse {
	return models.PersonResponse{
		ID:      p.Id,
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address,
		Work:    p.Work,
	}
}
