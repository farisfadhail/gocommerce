package migrations

import (
	"fmt"
	"gocommerce/database"
	"gocommerce/models/entity"
	"log"
)

func RunMigration() {
	db := database.DatabaseInit()

	err := db.AutoMigrate(&entity.User{}, &entity.Category{}, &entity.Product{})

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database Migrated")
}
