package migrations

import (
	"fmt"
	"gocommerce/database"
	"gocommerce/models/entity"
)

func RunMigration() {
	db := database.DatabaseInit()

	err := db.AutoMigrate(&entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ImageGallery{}, &entity.Cart{}, &entity.UserOrder{}, &entity.Order{}, &entity.Payment{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database Migrated")
}

func RunMigrationDrop() {
	db := database.DatabaseInit()

	err := db.Migrator().DropTable(&entity.User{}, &entity.Category{}, &entity.Product{}, &entity.ImageGallery{}, &entity.Cart{}, &entity.UserOrder{}, &entity.Order{}, &entity.Payment{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database Freshed")
}

func RunIndexES() {
	es := database.ElasticsearchInit()

	_, err := es.Indices.Create("gocommerce-products")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Index Created")
}

func RunDeleteIndexES() {
	es := database.ElasticsearchInit()

	_, err := es.Indices.Delete([]string{"gocommerce-products"})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Index Deleted")
}
