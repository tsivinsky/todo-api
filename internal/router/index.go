package router

import "github.com/gofiber/fiber/v2"

func IndexRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}
