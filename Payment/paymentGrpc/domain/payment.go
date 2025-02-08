package domain

type AlreadyPayment struct {
	Amount                    float32 `json:"amount"`
	CreditCardNumber          string  `json:"credit_card_number"`           // 信用卡号码
	CreditCardCVV             int32   `json:"credit_card_cvv"`              // CVV 码
	CreditCardExpirationYear  int32   `json:"credit_card_expiration_year"`  // 到期年份
	CreditCardExpirationMonth int32   `json:"credit_card_expiration_month"` // 到期月份
	OrderID                   string  `json:"order_id"`                     // 订单 ID
	TransactionID             string  `json:"transaction_id"`               // 交易 ID
	UserID                    uint32  `json:"user_id"`                      // 用户 ID
}

// TableName 指定 GORM 使用的表名
func (AlreadyPayment) TableName() string {
	return "already_payments" // 使用复数形式的表名
}
