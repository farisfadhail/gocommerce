package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gocommerce/models/entity"
	"gocommerce/models/request"
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

	var product entity.Product

	result := db.Debug().First(&product, "slug = ?", productSlug)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested Slug has been retrieved.",
		"data":    product,
	})
}

func UpdateProductHandler(ctx *fiber.Ctx) error {
	panic("udin")
}

func DeleteProductHandler(ctx *fiber.Ctx) error {
	panic("udin")
}
