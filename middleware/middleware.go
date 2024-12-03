package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/utils"
)

func Authenticated(ctx *fiber.Ctx) error {
	token := ctx.Get("x-auth-token")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated.",
		})
	}

	_, err := utils.VerifyTokenJwt(token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated.",
		})
	}

	claims, err := utils.DecodeToken(token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated.",
		})
	}

	if claims["exp"].(int64) >= utils.GetTimeNow() {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token has expired.",
		})
	}

	return ctx.Next()
}

func IsAdmin(ctx *fiber.Ctx) error {
	token := ctx.Get("x-auth-token")
	claims, err := utils.DecodeToken(token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated.",
		})
	}

	role := claims["role"].(string)
	if role != "admin" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden access.",
		})
	}

	return ctx.Next()
}
