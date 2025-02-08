package VO

// CreditCardInfo 用于接收信用卡信息
type CreditCardInfo struct {
	CreditCardNumber          string `json:"credit_card_number"`           // 信用卡号码
	CreditCardCVV             int32  `json:"credit_card_cvv"`              // CVV 码
	CreditCardExpirationYear  int32  `json:"credit_card_expiration_year"`  // 到期年份
	CreditCardExpirationMonth int32  `json:"credit_card_expiration_month"` // 到期月份
}

// ChargeReq 用于接收支付请求
type ChargeReq struct {
	Amount     float32        `json:"amount"`      // 支付金额
	CreditCard CreditCardInfo `json:"credit_card"` // 信用卡信息
	OrderID    string         `json:"order_id"`    // 订单 ID
	UserID     uint32         `json:"user_id"`     // 用户 ID
}

// ChargeResp 用于响应支付请求
type ChargeResp struct {
	TransactionID string `json:"transaction_id"` // 交易 ID
}
