package migrations

import (
	"fmt"
	"gocommerce/database"
	"gocommerce/models/entity"
)

func RunMigration() {
	db := database.DatabaseInit()

	//db.Migrator().DropTable(&entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ImageGallery{}, &entity.Cart{}, &entity.UserOrder{}, &entity.Order{}, &entity.Payment{})
	//fmt.Println("Database Freshed")

	err := db.AutoMigrate(&entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ImageGallery{}, &entity.Cart{}, &entity.UserOrder{}, &entity.Order{}, &entity.Payment{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database Migrated")
}
