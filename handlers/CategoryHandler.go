package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gocommerce/models/request"
	"gocommerce/repository"
)

type CategoryHandler struct {
	Repository repository.CategoryRepository
	Validate   *validator.Validate
}

func NewCategoryHandler(repo repository.CategoryRepository, validate *validator.Validate) *CategoryHandler {
	return &CategoryHandler{
		Repository: repo,
		Validate:   validate,
	}
}

func (c *CategoryHandler) GetAllCategoriesHandler(ctx *fiber.Ctx) error {
	categories, err := c.Repository.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve category data.",
			"error":   err.Error(),
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

func (c *CategoryHandler) StoreCategoryHandler(ctx *fiber.Ctx) error {
	categoryRequest := new(request.CategoryRequest)

	if err := ctx.BodyParser(categoryRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	if err := c.Validate.Struct(categoryRequest); err != nil {
		if validationError, ok := err.(*validator.InvalidValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid validation error.",
				"error":   FormatValidationError(validationError),
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unexpected validation error",
		})
	}

	newCategory, err := c.Repository.Create(*categoryRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to store category data.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New entry has been added to the database.",
		"data":    newCategory,
	})
}

func (c *CategoryHandler) GetBySlugCategoryHandler(ctx *fiber.Ctx) error {
	categorySlug := ctx.Params("categorySlug")

	category, err := c.Repository.GetBySlug(categorySlug)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve category data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested Slug has been retrieved.",
		"data":    category,
	})
}

func (c *CategoryHandler) UpdateCategoryHandler(ctx *fiber.Ctx) error {
	categoryRequest := new(request.CategoryUpdateRequest)

	if err := ctx.BodyParser(categoryRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	if err := c.Validate.Struct(categoryRequest); err != nil {
		if validationError, ok := err.(*validator.InvalidValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid validation error.",
				"error":   FormatValidationError(validationError),
			})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unexpected validation error",
		})
	}

	id := ctx.Params("categoryId")

	category, err := c.Repository.Update(*categoryRequest, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update category data.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes have been saved.",
		"data":    category,
	})
}

func (c *CategoryHandler) DeleteCategoryHandler(ctx *fiber.Ctx) error {
	categoryId := ctx.Params("categoryId")

	err := c.Repository.Delete(categoryId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete category data.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry has been removed from the database.",
	})
}
