package domain

import (
	"time"
)

type Cart struct {
	ID         int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int32     `gorm:"not null" json:"user_id"`
	ProductID  int32     `gorm:"not null" json:"product_id"`
	Quantity   int32     `gorm:"not null" json:"quantity"`
	Deleted    bool      `gorm:"default:false" json:"deleted"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
}
