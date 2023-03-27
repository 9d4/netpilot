package server

import (
	"fmt"
	"github.com/9d4/netpilot/database"
	"github.com/9d4/netpilot/ros/board"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	apiRouter.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	board.NewBoardHandler(database.DB()).SetupRoutes(apiRouter)

	app.Static("/", "./web/.output/public", fiber.Static{
		Compress: true,
		Index:    "index.html",
	})
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./web/.output/public/index.html", true)
	})
}

// checkError breaks the engine with log.Fatal if err != nil
func checkError(err error) {
	if err != nil {
		jww.ERROR.Fatal(err)
	}
}
