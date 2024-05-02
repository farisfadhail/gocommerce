package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gocommerce/models/entity"
	"gocommerce/models/request"
	"gocommerce/utils"
	"log"
	"strings"
	"time"
)

func RegisterHandler(ctx *fiber.Ctx) error {
	user := new(request.RegisterRequest)
	err := ctx.BodyParser(user)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	err = validate.Struct(user)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	newUser := entity.User{
		ID:       uuid.New(),
		FullName: user.FullName,
		Username: strings.ToLower(user.Username),
		Email:    user.Email,
		Role:     "consumer",
	}

	hashedPassword, err := utils.HashingPassword(user.Password)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error.",
		})
	}

	newUser.Password = hashedPassword

	result := db.Debug().Create(&newUser)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register user.",
		})
	}

	err = LoginHandler(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to login.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New user has been added to the database.",
		"data":    newUser,
	})
}

func LoginHandler(ctx *fiber.Ctx) error {
	loginRequest := new(request.LoginRequest)

	err := ctx.BodyParser(loginRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	err = validate.Struct(loginRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	// CHECK AVAILABLE USER
	var user entity.User
	result := db.First(&user, "email = ?", loginRequest.Email)

	if result.Error != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Wrong credential.",
		})
	}

	// CHECK VALIDATION PASSWORD
	isValid := utils.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Wrong credential.",
		})
	}

	// GENERATE JWT
	claims := jwt.MapClaims{}
	claims["name"] = user.FullName
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(60 * time.Minute).Unix()

	token, err := utils.GenerateJwtToken(&claims)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Wrong credential.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login success.",
		"token":   token,
	})
}

func UpdatePasswordUserHandler(ctx *fiber.Ctx) error {
	passwordRequest := new(request.UpdatePassword)
	err := ctx.BodyParser(passwordRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	var user entity.User

	userId := ctx.Params("userId")
	result := db.First(&user, userId)
	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	isValid := utils.CheckPasswordHash(passwordRequest.OldPassword, user.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Wrong password.",
		})
	}

	hashedPassword, err := utils.HashingPassword(passwordRequest.NewPassword)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error.",
		})
	}

	user.Password = hashedPassword
	result = db.Debug().Save(&user)
	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update password.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been updated.",
	})
}
