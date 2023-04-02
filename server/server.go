package server

import (
	"context"
	"fmt"
	"github.com/9d4/netpilot/database"
	"github.com/9d4/netpilot/ros/board"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	jww "github.com/spf13/jwalterweatherman"
	"log"
)

// Start ignites the engine
func Start(config *Config) {
	shared.Config = config

	app := fiber.New()
	applyRoutes(app)

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

	// api fallback, put at the end
	apiRouter.Get("/*", func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/boards/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		pubsub := database.RedisCli().Subscribe(context.Background(), "ch:board:"+c.Params("id")+":system/resource")
		defer pubsub.Close()

		c.SetCloseHandler(func(code int, text string) error {
			return nil
		})

		connClosed := make(chan bool, 1)
		go func() {
			for {
				if _, _, err := c.ReadMessage(); err != nil {
					connClosed <- true
					close(connClosed)
					break
				}
			}
		}()

		ch := pubsub.Channel()

	loop:
		for {
			select {
			case <-connClosed:
				c.Close()
				break loop
			case msg := <-ch:
				if err := c.WriteMessage(1, []byte(msg.Payload)); err != nil {
					fmt.Println(err)
					break loop
				}
			}
		}
	}))

	//staticRouter := app.Group("/")

	app.Route("/", func(router fiber.Router) {
		router.Static("/", "./web/dist", fiber.Static{
			Compress: true,
			Index:    "index.html",
		})
		router.Get("/*", func(c *fiber.Ctx) error {
			return c.SendFile("./web/dist/index.html", true)
		})

	})
}

// checkError breaks the engine with log.Fatal if err != nil
func checkError(err error) {
	if err != nil {
		jww.ERROR.Fatal(err)
	}
}
