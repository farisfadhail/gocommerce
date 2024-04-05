package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/utils"
)

func Authenticated(ctx *fiber.Ctx) error {
	token := ctx.Get("x-jwt-token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "UNAUTHENTICATED",
		})
	}

	_, err := utils.VerifyTokenJwt(token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "UNAUTHENTICATED",
		})
	}

	return ctx.Next()
}

func IsAdmin(ctx *fiber.Ctx) error {
	token := ctx.Get("x-jwt-token")
	claims, err := utils.DecodeToken(token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "UNAUTHENTICATED",
		})
	}

	role := claims["role"].(string)
	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "FORBIDDEN ACCESS",
		})
	}

	return ctx.Next()
}
