package types

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) handleGetUnits(c echo.Context) error {
	units, err := s.store.GetUnits(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, units)
}

func (s Server) handleGetNutrients(c echo.Context) error {
	nutrientTypes, err := s.store.GetNutrientTypes(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, nutrientTypes)
}

func (s Server) handleGetVitamins(c echo.Context) error {
	vitaminTypes, err := s.store.GetVitaminTypes(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, vitaminTypes)
}

func (s Server) handleGetMinerals(c echo.Context) error {
	mineralTypes, err := s.store.GetMineralTypes(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, mineralTypes)
}
