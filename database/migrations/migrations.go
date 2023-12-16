package migrations

import (
	"ecommerce-golang/database"
	"fmt"
	"log"
)

func RunMigration() {
	db := database.DatabaseInit()

	err := db.AutoMigrate()

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database Migrated")
}
