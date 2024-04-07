package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gocommerce/models/entity"
	"gocommerce/models/request"
)

func GetAllCartsHandler(ctx *fiber.Ctx) error {
	var carts []entity.Cart

	result := db.Debug().Find(&carts)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all data.",
		})
	}

	if len(carts) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Carts data is empty",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All cart data has been retrieved successfully.",
		"data":    carts,
	})
}

func StoreCartHandler(ctx *fiber.Ctx) error {
	cartRequest := new(request.CartRequest)
	err := ctx.BodyParser(cartRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(cartRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	var user entity.User
	result := db.First(&user, cartRequest.UserId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found.",
		})
	}

	var product entity.Product

	result = db.First(&product, cartRequest.ProductId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found.",
		})
	}

	var cart entity.Cart

	result = db.First(&cart, "product_id = ? AND user_id = ?", cartRequest.ProductId, cartRequest.UserId)

	if result.Error != nil {
		newCart := entity.Cart{
			UserId:    cartRequest.UserId,
			ProductId: cartRequest.ProductId,
			Quantity:  cartRequest.Quantity,
		}

		result = db.Debug().Create(&newCart)

		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create data.",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "New entry has been added to the database.",
			"data":    newCart,
		})
	}

	cart.Quantity += cartRequest.Quantity

	result = db.Debug().Save(&cart)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes has been save.",
		"data":    cart,
	})
}

func ShowByUserIdCartHandler(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	var carts []entity.Cart

	result := db.Debug().Find(&carts, "user_id = ?", userId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	if len(carts) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "You have no any products in your cart.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested user has been retrieved.",
		"data":    carts,
	})
}

func ShowByIdCartHandler(ctx *fiber.Ctx) error {
	cartId := ctx.Params("cartId")

	var cart entity.Cart

	result := db.Debug().First(&cart, cartId)

	var responseData []interface{}

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	responseData = append(responseData, cart)

	var product entity.Product

	result = db.Debug().First(&product, cart.ProductId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	responseData = append(responseData, product)

	var user entity.User

	result = db.Debug().First(&user, cart.UserId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	responseData = append(responseData, user)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data with the requested ID has been retrieved.",
		"data":    responseData,
	})
}

func UpdateQuantityCartHandler(ctx *fiber.Ctx) error {
	cartRequest := new(request.CartUpdateRequest)
	err := ctx.BodyParser(cartRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	// Validate request data
	err = validate.Struct(cartRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	cartId := ctx.Params("cartId")

	var cart entity.Cart

	result := db.First(&cart, cartId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	cart.Quantity = cartRequest.Quantity

	result = db.Debug().Save(&cart)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update cart data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes have been saved.",
		"data":    cart,
	})
}

func DeleteCartHandler(ctx *fiber.Ctx) error {
	cartRequest := new(request.CartDeleteRequest)
	err := ctx.BodyParser(cartRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(cartRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"data":    err.Error(),
		})
	}

	var cart entity.Cart

	for _, cartId := range cartRequest.CartId {
		result := db.First(&cart, cartId)

		if result.Error != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data not found.",
			})
		}
	}

	for _, cartId := range cartRequest.CartId {
		result := db.Debug().Delete(&cart, cartId)

		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to delete data.",
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry has been removed from the database.",
	})
}
