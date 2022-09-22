package handlers

import (
	"entetry/gotest/internal/dto"
	"entetry/gotest/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Company struct {
	companyService *service.Company
}

func NewCompany(companyService *service.Company) *Company {
	return &Company{companyService: companyService}
}

func (c *Company) GetById(ctx echo.Context) error {
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
	request := new(dto.AddCompanyRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.companyService.Create(ctx.Request().Context(), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, "Company created")
}

func (c *Company) Update(ctx echo.Context) error {
	request := new(dto.UpdateCompanyRequest)
	err := ctx.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = c.companyService.Update(ctx.Request().Context(), request)

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
