package paymentMapper

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"os"
	"payment-grpc/checkout"
	"paymentGrpc/domain"
	"paymentGrpc/payment"
	"project1/Common/config"
)

type PaymentMapper struct {
	DB             *gorm.DB
	checkoutClient checkout.CheckoutServiceClient
}

func NewPaymentMapper() (*PaymentMapper, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	addr := os.Getenv("CHECKOUT_ADDRESS")
	port := os.Getenv("CHECKOUT_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &PaymentMapper{DB: db, checkoutClient: checkout.NewCheckoutServiceClient(conn)}, nil
}

func (pm *PaymentMapper) Create(req *payment.ChargeReq, resp *checkout.GetPaymentItemResp) (*domain.AlreadyPayment, error) {
	getpayment := &domain.AlreadyPayment{
		Amount:                    req.Amount,
		CreditCardCVV:             req.CreditCard.CreditCardCvv,
		CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
		CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		CreditCardNumber:          req.CreditCard.CreditCardNumber,
		OrderID:                   req.OrderId,
		TransactionID:             resp.TransactionID,
		UserID:                    req.UserId,
	}
	result := pm.DB.Model(&domain.AlreadyPayment{}).Create(&getpayment)
	if result.Error != nil {
		return nil, fmt.Errorf("创建支付信息失败" + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("创建支付信息失败")
	}
	return getpayment, nil
}
