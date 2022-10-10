package app

import (
	"log"
	"os"
	"strings"
	"todo-app/internal/router"

	"github.com/gofiber/fiber/v2"
)

func Start(app *fiber.App) {
	app.Get("/", router.IndexRoute)

	port := getPort(":5000")
	log.Fatal(app.Listen(port))
}

func getPort(fallbackPort string) string {
	port := os.Getenv("PORT")

	if port == "" {
		port = fallbackPort
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return port
}
