package router

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/router/gitea"
	"git.trap.jp/toki/bot_converter/router/github"
)

func (h *Handlers) postConverterGitHub(c echo.Context) error {
	converter := getConverter(c)
	config := getConverterConfig(c)

	var secret string
	if converter.Secret.Valid {
		secret = converter.Secret.String
	}
	msg, err := github.MakeMessage(c, config, secret)
	if err != nil {
		return err
	}

	if len(msg) > 0 {
		go func() {
			if _, err := h.postMessage(converter.ChannelID, msg); err != nil {
				log.Printf("An error occurred while sending message: %v\n", err)
			}
		}()
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handlers) postConverterGitea(c echo.Context) error {
	converter := getConverter(c)

	var secret string
	if converter.Secret.Valid {
		secret = converter.Secret.String
	}
	msg, err := gitea.MakeMessage(c, secret)
	if err != nil {
		if errors.Is(err, gitea.ErrBadSignature) {
			return echo.NewHTTPError(http.StatusUnauthorized, "bad signature")
		} else {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}

	if len(msg) > 0 {
		go func() {
			if _, err := h.postMessage(converter.ChannelID, msg); err != nil {
				log.Printf("An error occurred while sending message: %v\n", err)
			}
		}()
	}

	return c.NoContent(http.StatusNoContent)
}
