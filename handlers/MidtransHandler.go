package handlers

import (
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/spf13/viper"
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

	// 1. Set you ServerKey with globally
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
			BillInfo1: "Payment for :",
			BillInfo2: customer.Phone,
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
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}

	if paymentMethod == "Credit Card" && tokenId != "" {
		paymentType = coreapi.PaymentTypeCreditCard
		chargeReq.CreditCard = &coreapi.CreditCardDetails{
			TokenID: tokenId,
		}
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

	// 3. Request to Midtrans using global config
	coreApiRes, _ := coreapi.ChargeTransaction(&chargeReq)

	return coreApiRes, nil
}

func NotificationMidtrans() error {
	panic("udin")
}
