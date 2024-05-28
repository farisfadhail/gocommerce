package main

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/database"
	"gocommerce/database/migrations"
	"gocommerce/routes"
	"log"
)

func main() {
	// Database Init
	database.DatabaseInit()
	//migrations.RunMigrationDrop()
	migrations.RunMigration()

	// ElasticSearch Init
	database.ElasticsearchInit()
	//migrations.RunDeleteIndexES()
	migrations.RunIndexES()

	app := fiber.New()

	routes.RouteInit(app)

	err := app.Listen("localhost:3000")
	if err != nil {
		log.Println("Failed to listen Go Fiber Server")
	}
}
