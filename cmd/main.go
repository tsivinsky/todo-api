package main

import (
	"log"
	"todo-app/internal/app"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	app.Start(fiber.New(fiber.Config{}))
}
