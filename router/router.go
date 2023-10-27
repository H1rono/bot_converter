package router

import (
	"context"
	"github.com/motoki317/sc"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sapphi-red/go-traq"

	"git.trap.jp/toki/bot_converter/repository"
)

type Handlers struct {
	repo        repository.Repository
	api         *traq.APIClient
	auth        context.Context
	accessToken string

	throttle *sc.Cache[postMessageArgs, *traq.Message]
}

func SetUp(c Config, e *echo.Echo, repo repository.Repository) {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, c.AccessToken)

	h := &Handlers{
		repo:        repo,
		api:         client,
		auth:        auth,
		accessToken: c.AccessToken,
	}

	h.throttle = sc.NewMust(h._postMessage, 5*time.Second, 5*time.Second, sc.WithCleanupInterval(1*time.Minute))

	retrieveCID := retrieveConverterID(h)

	e.GET("/", h.GetRoot)
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
