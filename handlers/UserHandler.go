package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gocommerce/database"
	"gocommerce/models/entity"
	"gocommerce/models/request"
	"gocommerce/models/response"
	"strings"
)

var db = database.DatabaseInit()
var validate = validator.New()

func GetAllUserHandler(ctx *fiber.Ctx) error {
	var users []entity.User

	result := db.Debug().Find(&users)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to Get All Datas",
		})
	}

	if len(users) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "All data is empty.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All data has been retrieved successfully.",
		"data":    users,
	})
}

func GetByIdUserHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	var user entity.User

	result := db.Debug().First(&user, userId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data Not Found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested ID has been retrieved.",
		"data":    user,
	})
}

func UpdateUserHandler(ctx *fiber.Ctx) error {
	userRequest := new(request.UserUpdateRequest)
	err := ctx.BodyParser(userRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
		})
	}

	var user entity.User

	userId := ctx.Params("userId")

	result := db.First(&user, userId)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	user.FullName = userRequest.FullName
	user.Username = strings.ToLower(userRequest.Username)

	result = db.Debug().Save(&user)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes have been saved.",
		"data":    user,
	})
}

func UpdateEmailUserHandler(ctx *fiber.Ctx) error {
	userRequest := new(request.UserUpdateEmailRequest)
	err := ctx.BodyParser(userRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
		})
	}

	var user entity.User

	userId := ctx.Params("userId")
	result := db.Debug().First(&user, userId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	if userRequest.Email != "" {
		err := validate.Struct(userRequest)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		user.Email = userRequest.Email
	}

	result = db.Debug().Save(&user)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update email",
		})
	}

	userResponse := response.UpdateEmailUserResponse{
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email data has been successfully updated.",
		"data":    userResponse,
	})
}

func UpdateRoleUserHandler(ctx *fiber.Ctx) error {
	userRequest := new(request.UserUpdateRoleRequest)
	err := ctx.BodyParser(userRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
		})
	}

	var user entity.User

	userId := ctx.Params("userId")

	result := db.First(&user, userId)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	user.Role = userRequest.Role

	result = db.Debug().Save(&user)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Change role have been saved.",
		"data":    user,
	})
}

func DeleteUserHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	var user entity.User

	result := db.Debug().First(&user)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	result = db.Debug().Delete(&user, userId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete data",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry has been removed from the database.",
	})
}
