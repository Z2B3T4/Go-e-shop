package domain

type CartItem struct {
	ProductID uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type CartVO struct {
	UserID uint32   `json:"user_id"`
	Item   CartItem `json:"item"`
}
