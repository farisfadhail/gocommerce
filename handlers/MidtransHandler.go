package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/spf13/viper"
	"gocommerce/models/entity"
	"strings"
)

func CreateChargeTransaction(paymentMethod string, orderId string, grossAmount int64, items []midtrans.ItemDetails, customer midtrans.CustomerDetails, customExpiry *coreapi.CustomExpiry, tokenId string) (*coreapi.ChargeResponse, error) {
	viper.SetConfigFile("config.env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var serverKey = viper.GetString("MIDTRANS_SERVER_KEY")

	midtrans.ServerKey = serverKey
	midtrans.Environment = midtrans.Sandbox

	var chargeReq coreapi.ChargeReq

	var paymentType coreapi.CoreapiPaymentType
	switch paymentMethod {
	case "Permata Virtual Account":
		paymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: "permata",
		}
	case "BCA Virtual Account":
		paymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: "bca",
		}
	case "Mandiri Bill Payment":
		paymentType = coreapi.PaymentTypeEChannel
		chargeReq.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "Payment for : ",
			BillInfo2: orderId,
		}
	case "BNI Virtual Account":
		paymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: "bni",
		}
	case "BRI Virtual Account":
		paymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: "bri",
		}
	case "BCA KlikPay":
		paymentType = coreapi.PaymentTypeBCAKlikpay
	case "KlikBCA":
		paymentType = coreapi.PaymentTypeKlikBca
	case "BRImo":
		paymentType = coreapi.PaymentTypeBRIEpay
	case "CIMB Clicks":
		paymentType = coreapi.PaymentTypeCimbClicks
	case "Danamon Online Banking":
		paymentType = coreapi.PaymentTypeDanamonOnline
	case "QRIS":
		paymentType = coreapi.PaymentTypeQris
	case "GoPay":
		paymentType = coreapi.PaymentTypeGopay
	case "ShopeePay":
		paymentType = coreapi.PaymentTypeShopeepay
	case "Indomaret", "Alfamart":
		paymentType = coreapi.PaymentTypeConvenienceStore
		chargeReq.ConvStore = &coreapi.ConvStoreDetails{
			Store: paymentMethod,
		}
	case "Akulaku":
		paymentType = coreapi.PaymentTypeAkulaku
	case "Credit Card":
		if tokenId == "" {
			return nil, fmt.Errorf("Token id is required for credit card payment method.")
		}

		paymentType = coreapi.PaymentTypeCreditCard
		chargeReq.CreditCard = &coreapi.CreditCardDetails{
			TokenID: tokenId,
		}
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}

	chargeReq.TransactionDetails = midtrans.TransactionDetails{
		OrderID:  orderId,
		GrossAmt: grossAmount,
	}

	chargeReq.Items = &items

	chargeReq.CustomerDetails = &customer

	if customExpiry != nil {
		chargeReq.CustomExpiry = customExpiry
	}

	chargeReq.PaymentType = paymentType

	coreApiRes, _ := coreapi.ChargeTransaction(&chargeReq)

	return coreApiRes, nil
}

func NotificationMidtrans(ctx *fiber.Ctx) error {
	// 1. Initialize empty map
	var notificationPayload map[string]interface{}

	// 2. Parse JSON request body and use it to set json to payload
	err := ctx.BodyParser(&notificationPayload)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request.",
			"error":   err.Error(),
		})

	}

	// 3. Get order-id from payload
	orderIdRaw, exists := notificationPayload["order_id"].(string)
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Not found.",
			"error":   err.Error(),
		})
	}

	if strings.Contains(orderIdRaw, ", ") {
		orderIds := strings.Split(orderIdRaw, ", ")
		for _, orderId := range orderIds {
			var order entity.Order

			db.First(&order, orderId)

			// 4. Check transaction to Midtrans with param orderId
			transactionStatusResp, err := coreapi.CheckTransaction(orderId)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error.",
					"error":   err.Error(),
				})
			}

			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
					order.Status = "Challenge"
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					order.Status = "Success"
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your database to 'success'
				order.Status = "Success"
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
				order.Status = "Pending"
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your database to 'failure'
				order.Status = "Failed"
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your database to 'pending' / waiting payment
				order.Status = "Pending"
			}

			result := db.Save(&order)

			if result.Error != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to update order status.",
					"error":   err.Error(),
				})
			}
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success update order status with midtrans webhook.",
		})
	}

	var order entity.Order

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, err := coreapi.CheckTransaction(orderIdRaw)
	if err != nil {
		//http.Error(ctx, err.Error(), http.StatusInternalServerError)
		return "", fmt.Errorf("Error checking transaction status")
	}

	if transactionStatusResp != nil {
		// 5. Do set transaction status based on response from check transaction status
		if transactionStatusResp.TransactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				order.Status = "Challenge"
			} else if transactionStatusResp.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				order.Status = "Success"
			}
		} else if transactionStatusResp.TransactionStatus == "settlement" {
			// TODO set transaction status on your database to 'success'
			order.Status = "Success"
		} else if transactionStatusResp.TransactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
			order.Status = "Pending"
		} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
			// TODO set transaction status on your database to 'failure'
			order.Status = "Failed"
		} else if transactionStatusResp.TransactionStatus == "pending" {
			// TODO set transaction status on your database to 'pending' / waiting payment
			order.Status = "Pending"
		}
	}

	result := db.Save(&order)

	if result.Error != nil {
		return "", fmt.Errorf("Failed to update order data.")
	}

	return "Success post Midtrans notification.", nil
}
