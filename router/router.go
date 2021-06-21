package router

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/repository"
)

type Handlers struct {
	Repo repository.Repository
}

func SetUp(e *echo.Echo, repo repository.Repository) {
	h := &Handlers{
		Repo: repo,
	}

	e.GET("/", h.GetRoot)
}

func (h *Handlers) GetRoot(c echo.Context) error {
	var res = struct {
		Message string `json:"message"`
	}{
		Message: "Hello, world!",
	}
	return c.JSON(http.StatusOK, res)
}
