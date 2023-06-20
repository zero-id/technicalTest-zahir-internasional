package handlers

import (
	contactsdto "contact/dto/contacts"
	dto "contact/dto/result"
	"contact/models"
	"contact/repositories"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	ContactRepository repositories.ContactRepository
}

func HandlerContact(ContactRepository repositories.ContactRepository) *handler {
	return &handler{ContactRepository}
}

func (h *handler) FindContacts(c echo.Context) error {

	page, _ := strconv.Atoi(c.QueryParam("page"))

	if page <= 0 {
		page = 1
	}

	perPage := 5

	offset := (page - 1) * perPage

	contacts, totalContacts, err := h.ContactRepository.FindContacts(offset, perPage)

	name := c.QueryParam("name")
	gender := c.QueryParam("gender")
	phone := c.QueryParam("phone")
	email := c.QueryParam("email")

	filteredContacts := []models.Contact{}
	for _, contact := range contacts {
		if name != "" && contact.Name != name {
			continue
		}
		if gender != "" && contact.Gender != gender {
			continue
		}
		if phone != "" && contact.Phone != phone {
			continue
		}
		if email != "" && contact.Email != email {
			continue
		}
		filteredContacts = append(filteredContacts, contact)
	}

	sortBy := c.QueryParam("sort")
	switch sortBy {
	case "name":
		// Mengurutkan kontak berdasarkan nama
		sort.SliceStable(filteredContacts, func(i, j int) bool {
			return filteredContacts[i].Name < filteredContacts[j].Name
		})
	default:
	}

	response := map[string]interface{}{
		"page":       page,
		"perPage":    perPage,
		"totalItems": len(filteredContacts),
		"totalPages": int(math.Ceil(float64(totalContacts) / float64(perPage))),
		"contacts":   filteredContacts,
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: response})
}

func (h *handler) GetContact(c echo.Context) error {
	id := c.Param("id")

	contact, err := h.ContactRepository.GetContact(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(contact)})
}

func convertResponse(u models.Contact) contactsdto.ContactResponse {
	return contactsdto.ContactResponse{
		ID:     u.ID,
		Name:   u.Name,
		Gender: u.Gender,
		Phone:  u.Phone,
		Email:  u.Email,
	}
}

func (h *handler) CreateContact(c echo.Context) error {
	request := new(contactsdto.CreateContactRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	contact := models.Contact{
		ID:     uuid.New().String(),
		Name:   request.Name,
		Gender: request.Gender,
		Phone:  request.Phone,
		Email:  request.Email,
	}

	data, err := h.ContactRepository.CreateContact(contact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)})
}

func (h *handler) UpdateContact(c echo.Context) error {
	request := contactsdto.UpdateContactRequest{
		Name:   c.FormValue("name"),
		Gender: c.FormValue("gender"),
		Phone:  c.FormValue("phone"),
		Email:  c.FormValue("email"),
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	id := c.Param("id")

	contact, err := h.ContactRepository.GetContact(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Name != "" {
		contact.Name = request.Name
	}

	if request.Gender != "" {
		contact.Gender = request.Gender
	}

	if request.Phone != "" {
		contact.Phone = request.Phone
	}

	if request.Email != "" {
		contact.Email = request.Email
	}

	data, err := h.ContactRepository.UpdateContact(contact)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)})
}

func (h *handler) DeleteContact(c echo.Context) error {
	id := c.Param("id")

	contact, err := h.ContactRepository.GetContact(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.ContactRepository.DeleteContact(contact, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(data)})
}
