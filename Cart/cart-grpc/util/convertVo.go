package util

import (
	"cart-grpc/domain"
	"time"
)

func ConvertToDb(cartVo *domain.CartVO) *domain.Cart {
	newcart := domain.Cart{
		UserID:     int32(cartVo.UserID),
		ProductID:  int32(cartVo.Item.ProductID),
		Quantity:   cartVo.Item.Quantity,
		Deleted:    false,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
	return &newcart
}
