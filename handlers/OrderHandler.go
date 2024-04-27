package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/spf13/cast"
	"gocommerce/models/entity"
	"gocommerce/models/request"
	"strconv"
	"strings"
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
	orderRequest := new(request.OrderRequest)
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

	var product entity.Product

	result := db.First(&product, orderRequest.ProductId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data product not found.",
		})
	}

	if uint64(orderRequest.Quantity) > product.Quantity {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Insufficient quantity for product.",
		})
	}

	var user entity.User

	result = db.First(&user, orderRequest.UserId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data user not found.",
		})
	}

	responseData := make(map[string]interface{})

	newUserOrder := entity.UserOrder{
		UserId:     orderRequest.UserId,
		Phone:      orderRequest.Phone,
		Address:    orderRequest.Address,
		District:   orderRequest.District,
		City:       orderRequest.City,
		Province:   orderRequest.Province,
		PostalCode: orderRequest.PostalCode,
	}

	result = db.Debug().Create(&newUserOrder)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user order data.",
		})
	}

	newOrder := entity.Order{
		ID:          uuid.New(),
		OrderNumber: "GCM" + newUserOrder.Phone + strconv.Itoa(int(time.Now().UnixMilli())),
		UserOrderId: newUserOrder.ID,
		ProductId:   orderRequest.ProductId,
		Quantity:    orderRequest.Quantity,
		Amount:      product.Price * int64(orderRequest.Quantity),
		Status:      "Pending",
	}

	result = db.Debug().Create(&newOrder)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create order data.",
		})
	}

	chargeResponse, err := CreateChargeTransaction(orderRequest.PaymentType, cast.ToString(newOrder.ID), int64(newOrder.Amount), []midtrans.ItemDetails{
		{
			Name:  product.Name,
			Price: product.Price,
			Qty:   int32(newOrder.Quantity),
		},
	}, midtrans.CustomerDetails{
		FName: user.FullName,
		Email: user.Email,
		Phone: newUserOrder.Phone,
		BillAddr: &midtrans.CustomerAddress{
			FName:    user.FullName,
			Phone:    newUserOrder.Phone,
			Address:  newUserOrder.Address,
			City:     newUserOrder.City,
			Postcode: strconv.Itoa(newUserOrder.PostalCode),
		},
		ShipAddr: &midtrans.CustomerAddress{
			FName:    user.FullName,
			Phone:    newUserOrder.Phone,
			Address:  newUserOrder.Address,
			City:     newUserOrder.City,
			Postcode: strconv.Itoa(newUserOrder.PostalCode),
		},
	}, nil, orderRequest.TokenID)

	if err != nil {
		db.Delete(&entity.UserOrder{}, newUserOrder.ID)

		db.Delete(&entity.Order{}, newOrder.ID)

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create charge transaction.",
		})
	}

	newPayment := entity.Payment{
		ID:              uuid.New(),
		TransactionID:   chargeResponse.TransactionID,
		Amount:          newOrder.Amount,
		PaymentType:     chargeResponse.PaymentType,
		Status:          chargeResponse.TransactionStatus,
		TransactionTime: chargeResponse.TransactionTime,
	}

	result = db.Debug().Create(&newPayment)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create data.",
		})
	}

	product.Quantity = product.Quantity - uint64(orderRequest.Quantity)
	result = db.Debug().Save(&product)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update product data.",
		})
	}

	responseData["user_order"] = newUserOrder
	responseData["order"] = newOrder
	responseData["payment"] = chargeResponse

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order has been created successfully. Finish payment to process your order.",
		"data":    responseData,
	})
}

// Checkout by cart
func CheckoutByCartOrderHandler(ctx *fiber.Ctx) error {
	orderRequest := new(request.OrderByCartRequest)
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

	var user entity.User

	result := db.First(&user, orderRequest.UserId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data user not found.",
		})
	}

	responseData := make(map[string]interface{})

	var cart []entity.Cart

	result = db.Find(&cart, orderRequest.CartId)

	if result.Error != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Cart not found",
		})
	}

	// Check if product quantity is enough
	for _, cartDetail := range cart {
		var product entity.Product

		db.First(&product, cartDetail.ProductId)

		if product.Quantity < uint64(cartDetail.Quantity) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Insufficient quantity for product.",
			})
		}
	}

	var Amount int64

	var items []midtrans.ItemDetails

	newUserOrder := entity.UserOrder{
		UserId:     orderRequest.UserId,
		Phone:      orderRequest.Phone,
		Address:    orderRequest.Address,
		District:   orderRequest.District,
		City:       orderRequest.City,
		Province:   orderRequest.Province,
		PostalCode: orderRequest.PostalCode,
	}

	var orderIDS []string

	for _, cartDetail := range cart {
		var product entity.Product

		db.First(&product, cartDetail.ProductId)

		newOrder := entity.Order{
			ID:          uuid.New(),
			OrderNumber: "GCM" + newUserOrder.Phone + strconv.Itoa(int(time.Now().UnixMilli())),
			UserOrderId: newUserOrder.UserId,
			ProductId:   cartDetail.ProductId,
			Quantity:    cartDetail.Quantity,
			Amount:      product.Price * int64(cartDetail.Quantity),
			Status:      "Pending",
		}

		result := db.Debug().Create(&newOrder)

		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create order data.",
			})
		}

		orderIDS = append(orderIDS, cast.ToString(newOrder.ID))

		Amount = Amount + (product.Price * int64(cartDetail.Quantity))
		items = append(items, midtrans.ItemDetails{
			Name:  product.Name,
			Price: product.Price,
			Qty:   int32(cartDetail.Quantity),
		})

		product.Quantity = product.Quantity - uint64(cartDetail.Quantity)
		db.Save(&product)
	}

	result = db.Debug().Create(&newUserOrder)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user order data.",
		})
	}

	chargeResponse, err := CreateChargeTransaction(orderRequest.PaymentType, strings.Join(orderIDS, ", "), Amount, items, midtrans.CustomerDetails{
		FName: user.FullName,
		Email: user.Email,
		Phone: newUserOrder.Phone,
		BillAddr: &midtrans.CustomerAddress{
			FName:    user.FullName,
			Phone:    newUserOrder.Phone,
			Address:  newUserOrder.Address,
			City:     newUserOrder.City,
			Postcode: strconv.Itoa(newUserOrder.PostalCode),
		},
		ShipAddr: &midtrans.CustomerAddress{
			FName:    user.FullName,
			Phone:    newUserOrder.Phone,
			Address:  newUserOrder.Address,
			City:     newUserOrder.City,
			Postcode: strconv.Itoa(newUserOrder.PostalCode),
		},
	},
		nil, orderRequest.TokenID)

	if err != nil {
		db.Delete(&entity.UserOrder{}, newUserOrder.ID)

		for _, orderID := range orderIDS {
			db.Delete(&entity.Order{}, orderID)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create charge transaction.",
		})
	}

	newPayment := entity.Payment{
		ID:              uuid.New(),
		TransactionID:   chargeResponse.TransactionID,
		Amount:          Amount,
		PaymentType:     chargeResponse.PaymentType,
		Status:          "Pending",
		TransactionTime: chargeResponse.TransactionTime,
	}

	result = db.Debug().Create(&newPayment)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create data.",
		})
	}

	responseData["user_order"] = newUserOrder
	responseData["order"] = cart
	responseData["payment"] = chargeResponse

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New entry has been added to the database.",
		"data":    responseData,
	})
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
