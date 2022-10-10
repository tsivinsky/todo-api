package main

import (
	"todo-app/internal/app"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app.Start(fiber.New(fiber.Config{}))
}
