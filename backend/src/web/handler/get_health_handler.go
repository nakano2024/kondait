package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetHealthHandler struct{}

func NewGetHealthHandler() *GetHealthHandler {
	return &GetHealthHandler{}
}

func (h *GetHealthHandler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
