package router

import (
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/model"
	"git.trap.jp/toki/bot_converter/repository"
)

const (
	converterKey       = "converter"
	converterConfigKey = "converter_config"
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
				if errors.Is(err, repository.ErrNotFound) {
					return echo.NewHTTPError(http.StatusNotFound, "converter not found")
				} else {
					c.Logger().Error(err)
					return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
				}
			}

			config, err := h.repo.GetConverterConfig(cid)
			if err != nil && !errors.Is(err, repository.ErrNotFound) {
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
			if errors.Is(err, repository.ErrNotFound) {
				config = &model.Config{ConverterID: cid}
			}

			c.Set(converterKey, converter)
			c.Set(converterConfigKey, config)

			return next(c)
		}
	}
}

func getConverter(c echo.Context) *model.Converter {
	return c.Get(converterKey).(*model.Converter)
}

func getConverterConfig(c echo.Context) *model.Config {
	return c.Get(converterConfigKey).(*model.Config)
}
