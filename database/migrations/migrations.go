package migrations

import (
	"fmt"
	"gocommerce/database"
	"gocommerce/models/entity"
	"log"
)

func RunMigration() {
	db := database.DatabaseInit()

	//db.Migrator().DropTable(&entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ImageGallery{}, &entity.Cart{})
	//fmt.Println("Database Freshed")

	err := db.AutoMigrate(&entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ImageGallery{}, &entity.Cart{})

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Database Migrated")
}
