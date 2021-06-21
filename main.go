package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/repository"
	"git.trap.jp/toki/bot_converter/router"
	"git.trap.jp/toki/bot_converter/service"
)

func main() {
	// db
	db, err := initDB()
	if err != nil {
		log.Fatalf("an error occurred while initializing db: %s", err)
	}
	repo := repository.NewGormRepository(db)

	// router
	e := echo.New()
	botHandler := service.SetUp(provideBotConfig(), repo)
	router.SetUp(provideRouterConfig(), e, repo, botHandler)

	if err := e.Start(fmt.Sprintf(":%d", c.Port)); err != nil {
		log.Fatalf("an error occurred while starting server: %s", err)
	}
}
