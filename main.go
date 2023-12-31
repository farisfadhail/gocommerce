package main

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/routes"
	"log"
)

func main() {
	app := fiber.New()

	routes.RouteInit(app)

	err := app.Listen("localhost:3000")
	if err != nil {
		log.Println("Failed to listen Go Fiber Server")
	}
}
