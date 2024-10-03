package products

import (
	"errors"
	v1 "github.com/kobietka/product-service/pkg/api/v1"
	"github.com/kobietka/product-service/pkg/text"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	searchLimit int8 = 15
)

type productBinding struct {
	Ean string `param:"ean"`
}

type searchBinding struct {
	Query string `query:"query"`
	Limit int8   `query:"limit"`
}

func (s Server) handleGetProduct(c echo.Context) error {
	var binding productBinding
	if err := c.Bind(&binding); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	product, err := s.store.GetProduct(c.Request().Context(), binding.Ean)
	if err != nil {
		if errors.Is(err, v1.ErrorDataNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, product)
}

func (s Server) handleSearchProduct(c echo.Context) error {
	var binding searchBinding
	if err := c.Bind(&binding); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if text.IsBlankString(binding.Query) {
		return c.NoContent(http.StatusBadRequest)
	}

	if binding.Limit <= 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	limit := min(binding.Limit, searchLimit)

	products, err := s.store.SearchProducts(c.Request().Context(), binding.Query, limit)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, products)
}

func (s Server) handlePostProduct(c echo.Context) error {
	var product v1.Product
	if err := c.Bind(&product); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := validateProduct(product); err != nil {
		return c.JSON(http.StatusBadRequest, v1.ErrorResponse{Code: err.Error()})
	}

	if err := s.store.CreateProduct(c.Request().Context(), product); err != nil {
		if errors.Is(err, v1.ErrorInvalidData) {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func (s Server) handlePutProduct(c echo.Context) error {
	var product v1.Product
	if err := c.Bind(&product); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := validateProduct(product); err != nil {
		return c.JSON(http.StatusBadRequest, v1.ErrorResponse{Code: err.Error()})
	}

	if err := s.store.UpdateProduct(c.Request().Context(), product); err != nil {
		if errors.Is(err, v1.ErrorInvalidData) {
			return c.NoContent(http.StatusBadRequest)
		}
		if errors.Is(err, v1.ErrorProductDoesNotExist) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s Server) handleDeleteProduct(c echo.Context) error {
	var binding productBinding
	if err := c.Bind(&binding); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := s.store.DeleteProduct(c.Request().Context(), binding.Ean); err != nil {
		if errors.Is(err, v1.ErrorProductDoesNotExist) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
