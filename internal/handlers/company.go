package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"entetry/gotest/internal/model"
	"entetry/gotest/internal/service"
)

// Company handler company struct
type Company struct {
	companyService *service.Company
}

// NewCompany creates new company handler
func NewCompany(companyService *service.Company) *Company {
	return &Company{companyService: companyService}
}

// GetAll return all companies
func (c *Company) GetAll(ctx echo.Context) error {
	companies, err := c.companyService.GetAll(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, companies)
}

// GetByID return by ID
func (c *Company) GetByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	company, err := c.companyService.GetByID(ctx.Request().Context(), id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, company)
}

// Create add company
func (c *Company) Create(ctx echo.Context) error {
	request := new(addCompanyRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	company := &model.Company{Name: request.Name}
	id, err := c.companyService.Create(ctx.Request().Context(), company)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, id)
}

// Update company
func (c *Company) Update(ctx echo.Context) error {
	request := new(updateCompanyRequest)
	err := ctx.Bind(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctx.Validate(request)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	company := &model.Company{ID: request.UUID, Name: request.Name}
	err = c.companyService.Update(ctx.Request().Context(), company)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, "Company updated")
}

// Delete Company
func (c *Company) Delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.companyService.Delete(ctx.Request().Context(), id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return ctx.JSON(http.StatusOK, "Company deleted")
}

// GetLogoByCompanyID get logo
func (c *Company) GetLogoByCompanyID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	logo, err := c.companyService.GetLogo(ctx.Request().Context(), id)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if logo == "" {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return ctx.File(logo)
}

// AddLogo add new logo
func (c *Company) AddLogo(ctx echo.Context) error {
	companyID := ctx.FormValue("companyID")
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	err = c.companyService.AddLogo(ctx.Request().Context(), companyID, file)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Logo has been added")
}
