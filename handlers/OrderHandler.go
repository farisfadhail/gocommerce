package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gocommerce/models/entity"
	"gocommerce/models/request"
	"strconv"
	"time"
)

func GetAllOrderHandler(ctx *fiber.Ctx) error {
	var orders []entity.Order

	result := db.Debug().Find(&orders)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all data.",
		})
	}

	if len(orders) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Order data is empty.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All order data has been retrieved successfully.",
		"data":    orders,
	})
}

// Checkout immediately
func CheckoutImmediatelyOrderHandler(ctx *fiber.Ctx) error {
	userOrderRequest := new(request.UserOrderRequest)
	err := ctx.BodyParser(userOrderRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(userOrderRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"data":    err.Error(),
		})
	}

	orderRequest := new(request.OrderRequest)
	err = ctx.BodyParser(orderRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(orderRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	var product entity.Product

	result := db.First(&product, orderRequest.ProductId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data product not found.",
		})
	}

	var user entity.User

	result = db.First(&user, userOrderRequest.UserId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data user not found.",
		})
	}

	var responseData []interface{}

	newUserOrder := entity.UserOrder{
		UserId:     userOrderRequest.UserId,
		Phone:      userOrderRequest.Phone,
		Address:    userOrderRequest.Address,
		District:   userOrderRequest.District,
		City:       userOrderRequest.City,
		Province:   userOrderRequest.Province,
		PostalCode: userOrderRequest.PostalCode,
	}

	result = db.Debug().Create(&newUserOrder)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user order data.",
		})
	}

	responseData = append(responseData, newUserOrder)

	newOrder := entity.Order{
		ID:          uuid.New(),
		OrderNumber: "GOCOM" + strconv.Itoa(int(newUserOrder.Phone)) + strconv.Itoa(int(time.Now().UnixMilli())),
		UserOrderId: newUserOrder.ID,
		ProductId:   orderRequest.ProductId,
		Quantity:    orderRequest.Quantity,
		TotalPrice:  product.Price * uint64(orderRequest.Quantity),
		Status:      "Pending",
	}

	result = db.Debug().Create(&newOrder)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create order data.",
		})
	}

	responseData = append(responseData, newOrder)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New entry has been added to the database.",
		"data":    responseData,
	})
}

// Checkout by cart
func CheckoutByCartOrderHandler(ctx *fiber.Ctx) error {
	panic("udin")
}

// update status if not accept yet
func UpdateOrderHandler(ctx *fiber.Ctx) error {
	orderRequest := new(request.OrderUpdateRequest)
	err := ctx.BodyParser(orderRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request is empty.",
		})
	}

	err = validate.Struct(orderRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})
	}

	orderNumber := ctx.Params("orderNumber")

	var order entity.Order

	result := db.First(&order, "order_number = ?", orderNumber)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data not found.",
		})
	}

	if order.Status == "Delivered" || order.Status == "Cancel" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Can't change order with status " + order.Status,
		})
	}

	order.Status = orderRequest.Status

	result = db.Debug().Save(&order)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update data.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Changes have been saved.",
		"data":    order,
	})
}
