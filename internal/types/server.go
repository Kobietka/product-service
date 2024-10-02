package types

import "github.com/labstack/echo/v4"

type Server struct {
	store Store
}

func NewServer(store Store) Server {
	return Server{store: store}
}

func (s Server) Routes(e *echo.Echo) {
	e.GET("/types/unit", s.handleGetUnits)
	e.GET("/types/nutrient", s.handleGetNutrients)
	e.GET("/types/vitamin", s.handleGetVitamins)
	e.GET("/types/mineral", s.handleGetMinerals)
}
