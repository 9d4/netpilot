package server

import (
	"fmt"
	"github.com/9d4/netpilot/database"
	"github.com/9d4/netpilot/ros/board"
	"github.com/gofiber/fiber/v2"
	jww "github.com/spf13/jwalterweatherman"
)

// Start ignites the engine
func Start(config *Config) {
	shared.Config = config

	app := fiber.New()
	applyRoutes(app)

	fmt.Println(config.v.GetString("ADDRESS"))
	checkError(app.Listen(config.v.GetString("ADDRESS")))
}

// define routes here
func applyRoutes(app *fiber.App) {
	apiRouter := app.Group("/api")

	board.NewBoardHandler(database.DB()).SetupRoutes(apiRouter)
}

// checkError breaks the engine with log.Fatal if err != nil
func checkError(err error) {
	if err != nil {
		jww.ERROR.Fatal(err)
	}
}
