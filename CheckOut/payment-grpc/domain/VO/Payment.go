package VO

// Address 表示地址信息
type Address struct {
	StreetAddress string `json:"street_address"` // 街道地址
	City          string `json:"city"`           // 城市
	State         string `json:"state"`          // 州
	Country       string `json:"country"`        // 国家
	ZipCode       string `json:"zip_code"`       // 邮政编码
}

// CreditCardInfo 表示信用卡信息
type CreditCardInfo struct {
	CreditCardNumber          string `json:"credit_card_number"`           // 信用卡号码
	CreditCardCVV             int32  `json:"credit_card_cvv"`              // CVV 码
	CreditCardExpirationYear  int32  `json:"credit_card_expiration_year"`  // 到期年份
	CreditCardExpirationMonth int32  `json:"credit_card_expiration_month"` // 到期月份
}

// CheckoutReq 表示结账请求
type CheckoutReq struct {
	UserID     uint32         `json:"user_id"`     // 用户 ID
	FirstName  string         `json:"firstname"`   // 名
	LastName   string         `json:"lastname"`    // 姓
	Email      string         `json:"email"`       // 电子邮件
	Address    Address        `json:"address"`     // 地址信息
	CreditCard CreditCardInfo `json:"credit_card"` // 信用卡信息
}

// CheckoutResp 表示结账响应
type CheckoutResp struct {
	OrderID       string `json:"order_id"`       // 订单 ID
	TransactionID string `json:"transaction_id"` // 交易 ID
}
