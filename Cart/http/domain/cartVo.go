package domain

type CartItem struct {
	ProductID uint32 `json:"productId"`
	Quantity  int32  `json:"quantity"`
}

type CartVO struct {
	UserID uint32   `json:"userId"`
	Item   CartItem `json:"item"`
}
