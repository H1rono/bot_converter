package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"git.trap.jp/toki/bot_converter/repository"
	"git.trap.jp/toki/bot_converter/router"
	"git.trap.jp/toki/bot_converter/service/bot"
)

func main() {
	// db
	db, err := initDB()
	if err != nil {
		log.Fatalf("an error occurred while initializing db: %s", err)
	}
	repo := repository.NewGormRepository(db)

	// bot
	go func() {
		if err := bot.Start(provideBotConfig(), repo); err != nil {
			panic(err)
		}
	}()

	// router
	e := echo.New()
	router.SetUp(provideRouterConfig(), e, repo)

	if err := e.Start(fmt.Sprintf(":%d", c.Port)); err != nil {
		log.Fatalf("an error occurred while starting server: %s", err)
	}
}
