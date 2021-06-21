package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sapphi-red/go-traq"

	"git.trap.jp/toki/bot_converter/repository"
)

type Handlers struct {
	repo        repository.Repository
	api         *traq.APIClient
	auth        context.Context
	accessToken string
}

func SetUp(c Config, e *echo.Echo, repo repository.Repository, botHandler echo.HandlerFunc) {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, c.AccessToken)

	h := &Handlers{
		repo:        repo,
		api:         client,
		auth:        auth,
		accessToken: c.AccessToken,
	}

	retrieveCID := retrieveConverterID(h)

	e.GET("/", h.GetRoot)
	e.POST("/bot", botHandler)
	convertersAPI := e.Group("/converters")
	{
		convertersCID := convertersAPI.Group("/:converterID", retrieveCID)
		{
			convertersCID.POST("/github", h.postConverterGitHub)
			convertersCID.POST("/gitea", h.postConverterGitea)
		}
	}
}

func (h *Handlers) GetRoot(c echo.Context) error {
	var res = struct {
		Message string `json:"message"`
	}{
		Message: "Hello, world!",
	}
	return c.JSON(http.StatusOK, res)
}
