package routes

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/config"
)

func RouteInit(app *fiber.App) {
	app.Static("/public", config.ProjectRootPath+"/public/asset")

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to ??????????")
	})
}
