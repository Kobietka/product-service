package products

import (
	"github.com/labstack/echo/v4"
)

type Server struct {
	store Store
}

func NewServer(store Store) Server {
	return Server{store: store}
}

func (s Server) Routes(e *echo.Echo) {
	e.GET("/products/:ean", s.handleGetProduct)
	e.GET("/products", s.handleSearchProduct)
	e.POST("/products", s.handlePostProduct)
	e.PUT("/products", s.handlePutProduct)
	e.DELETE("/products/:ean", s.handleDeleteProduct)
}
