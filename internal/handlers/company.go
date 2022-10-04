package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"entetry/gotest/internal/model"
	"entetry/gotest/internal/service"
)

type Company struct {
	companyService *service.Company
}

func NewCompany(companyService *service.Company) *Company {
	return &Company{companyService: companyService}
}

func (c *Company) GetAll(ctx echo.Context) error {
	companies, err := c.companyService.GetAll(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, companies)
}

func (c *Company) GetByID(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	company, err := c.companyService.GetById(ctx.Request().Context(), id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, company)
}

func (c *Company) Create(ctx echo.Context) error {
	request := new(AddCompanyRequest)
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

func (c *Company) Update(ctx echo.Context) error {
	request := new(UpdateCompanyRequest)
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
