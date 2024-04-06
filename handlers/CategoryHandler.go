package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gocommerce/models/entity"
	"gocommerce/models/request"
)

func GetAllCategoriesHandler(ctx *fiber.Ctx) error {
	var categories []entity.Category

	result := db.Debug().Find(&categories)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all data.",
		})
	}

	if len(categories) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Categories data is empty.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All category data has been retrieved successfully.",
		"data":    categories,
	})
}

func StoreCategoryHandler(ctx *fiber.Ctx) error {
	categoryRequest := new(request.CategoryRequest)
	err := ctx.BodyParser(categoryRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	err = validate.Struct(categoryRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	newCategory := entity.Category{
		Name: categoryRequest.Name,
		Slug: slug.Make(categoryRequest.Name),
	}

	result := db.Debug().Create(&newCategory)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to store category data",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New entry has been added to the database.",
		"data":    newCategory,
	})
}

func GetBySlugCategoryHandler(ctx *fiber.Ctx) error {
	categorySlug := ctx.Params("categorySlug")

	var category entity.Category

	result := db.Debug().First(&category, "slug = ?", categorySlug)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested Slug has been retrieved.",
		"data":    category,
	})
}

func UpdateCategoryHandler(ctx *fiber.Ctx) error {
	categoryRequest := new(request.CategoryUpdateRequest)
	err := ctx.BodyParser(categoryRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(categoryRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	categoryId := ctx.Params("categoryId")

	var category entity.Category

	result := db.First(&category, categoryId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	category.Name = categoryRequest.Name
	category.Slug = slug.Make(categoryRequest.Name)

	result = db.Debug().Save(&category)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update category data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes have been saved.",
		"data":    category,
	})
}

func DeleteCategoryHandler(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("categoryId")

	var category entity.Category

	result := db.First(&category, categoryId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	result = db.Debug().Delete(&category, categoryId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry has been removed from the database.",
	})
}
