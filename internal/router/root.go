package router

import "github.com/gofiber/fiber/v2"

func Init(app *fiber.App) {
	r := app.Group("/api")

	r.Get("/", handleIndex)

	r.Post("/auth/register", handleAuthRegister)
	r.Post("/auth/login", handleAuthLogin)
}

func handleIndex(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}
