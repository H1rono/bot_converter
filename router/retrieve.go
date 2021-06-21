package router

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/repository"
)

const (
	converterKey = "converter"
)

func retrieveConverterID(h *Handlers) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cid, err := uuid.FromString(c.Param("converterID"))
			if err != nil || cid == uuid.Nil {
				return echo.NewHTTPError(http.StatusBadRequest, "bad converter id")
			}

			converter, err := h.repo.GetConverter(cid)
			if err != nil {
				if err == repository.ErrNotFound {
					return echo.NewHTTPError(http.StatusNotFound, "converter not found")
				} else {
					c.Logger().Error(err)
					return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
				}
			}

			c.Set(converterKey, converter)

			return next(c)
		}
	}
}
