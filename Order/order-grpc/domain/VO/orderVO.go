package VO

// Address represents the address information of the user.
type Address struct {
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       int32  `json:"zip_code"`
}

// CartItem represents the item that is in the shopping cart.
type CartItem struct {
	ProductID uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

// OrderItem represents an item in the order.
type OrderItem struct {
	Item CartItem `json:"item"`
	Cost float32  `json:"cost"`
}

// PlaceOrderReq represents the request structure for placing an order.
type PlaceOrderReq struct {
	UserID       uint32      `json:"user_id"`
	UserCurrency string      `json:"user_currency"`
	Address      Address     `json:"address"`
	Email        string      `json:"email"`
	OrderItems   []OrderItem `json:"order_items"`
}

// OrderResult represents the result of placing an order.
type OrderResult struct {
	OrderID string `json:"order_id"`
}

// PlaceOrderResp represents the response structure for placing an order.
type PlaceOrderResp struct {
	Order OrderResult `json:"order"`
}
