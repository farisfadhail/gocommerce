package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gocommerce/models/entity"
	"gocommerce/models/request"
	"gocommerce/utils"
	"log"
)

func GetAllProductsHandler(ctx *fiber.Ctx) error {
	var products []entity.Product

	result := db.Debug().Find(&products)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all data.",
		})
	}

	if len(products) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Products data is empty.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All category data has been retrieved successfully.",
		"data":    products,
	})
}

func StoreProductHandler(ctx *fiber.Ctx) error {
	productRequest := new(request.ProductRequest)
	err := ctx.BodyParser(productRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(productRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	// validasi category id
	var category entity.Category

	result := db.First(&category, productRequest.CategoryId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   "Category is not found.",
		})
	}

	var responseData []interface{}

	newProduct := entity.Product{
		CategoryId:  productRequest.CategoryId,
		Name:        productRequest.Name,
		Slug:        slug.Make(productRequest.Name),
		Price:       productRequest.Price,
		Description: productRequest.Description,
		Quantity:    productRequest.Quantity,
	}

	result = db.Debug().Create(&newProduct)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create data.",
		})
	}

	responseData = append(responseData, newProduct)

	filenames := ctx.Locals("filenames").([]string)
	if len(filenames) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product images is required",
		})
	}

	for _, filename := range filenames {
		newImageGallery := entity.ImageGallery{
			ProductId: newProduct.ID,
			FileName:  filename,
		}

		result := db.Debug().Create(&newImageGallery)
		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to store image",
			})
		}

		responseData = append(responseData, newImageGallery)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New entry has been added to the database.",
		"data":    responseData,
	})
}

func GetBySlugProductHandler(ctx *fiber.Ctx) error {
	productSlug := ctx.Params("productSlug")

	var responseData []interface{}

	var product entity.Product

	result := db.Debug().First(&product, "slug = ?", productSlug)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	responseData = append(responseData, product)

	var imageGallery []entity.ImageGallery

	result = db.Debug().Find(&imageGallery, "product_id = ?", product.ID)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Can't get image gallery data.",
		})
	}

	for _, gallery := range imageGallery {
		responseData = append(responseData, gallery)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested Slug has been retrieved.",
		"data":    responseData,
	})
}

func UpdateProductHandler(ctx *fiber.Ctx) error {
	productRequest := new(request.ProductUpdateRequest)
	err := ctx.BodyParser(productRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(productRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	productId := ctx.Params("productId")

	var product entity.Product

	result := db.First(&product, productId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	product.CategoryId = productRequest.CategoryId
	product.Name = productRequest.Name
	product.Slug = slug.Make(productRequest.Name)
	product.Price = productRequest.Price
	product.Description = productRequest.Description
	product.Quantity = productRequest.Quantity

	result = db.Debug().Save(&product)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update product data.",
		})
	}

	var responseData []interface{}

	responseData = append(responseData, product)

	filenames := ctx.Locals("filenames").([]string)

	var imageGallery []entity.ImageGallery

	result = db.Debug().Find(&imageGallery, "product_id = ?", productId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Can't get image gallery data.",
		})
	}

	for _, gallery := range imageGallery {
		err := utils.HandleRemoveFile(gallery.FileName)

		if err != nil {
			log.Println("Failed delete image in directory.")
		}

		result = db.Debug().Delete(&gallery, "product_id = ?", productId)

		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Can't delete image gallery data.",
			})
		}
	}

	for _, filename := range filenames {
		newImageGallery := entity.ImageGallery{
			ProductId: product.ID,
			FileName:  filename,
		}

		result := db.Debug().Create(&newImageGallery)
		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to store image",
			})
		}

		responseData = append(responseData, newImageGallery)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes have been saved.",
		"data":    responseData,
	})
}

func DeleteProductHandler(ctx *fiber.Ctx) error {
	productId := ctx.Params("productId")

	var product entity.Product

	result := db.First(&product, productId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	var imageGallery []entity.ImageGallery

	result = db.Debug().Find(&imageGallery, "product_id = ?", productId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Can't get image gallery data.",
		})
	}

	for _, gallery := range imageGallery {
		err := utils.HandleRemoveFile(gallery.FileName)

		if err != nil {
			log.Println("Failed delete image in directory.")
		}

		result = db.Debug().Delete(&gallery, "product_id = ?", productId)

		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Can't delete image gallery data.",
			})
		}
	}

	result = db.Debug().Delete(&product, productId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry has been removed from the database.",
	})
}
