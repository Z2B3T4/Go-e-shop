package Entity

import "time"

// AddressEntity represents the database entity for addresses.
type AddressEntity struct {
	ID            int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int32     `json:"user_id"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       int32     `json:"zip_code"`
	Deleted       bool      `gorm:"default:false" json:"deleted"` // 逻辑删除字段
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderEntity represents the database entity for orders.
type OrderEntity struct {
	ID           int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int32     `json:"user_id"`
	AddressID    int32     `json:"address_id"`
	UserCurrency string    `json:"user_currency"` // 保留用户货币字段
	Email        string    `json:"email"`
	Deleted      bool      `gorm:"default:false" json:"deleted"` // 逻辑删除字段
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderItemEntity represents the database entity for order items.
type OrderItemEntity struct {
	ID        int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   int32     `json:"order_id"`
	ProductID int32     `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	Cost      float32   `json:"cost"`
	Deleted   bool      `gorm:"default:false" json:"deleted"` // 逻辑删除字段
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName sets the table name to be used for AddressEntity
func (AddressEntity) TableName() string {
	return "addresses"
}

// TableName sets the table name to be used for OrderEntity
func (OrderEntity) TableName() string {
	return "orders"
}

// TableName sets the table name to be used for OrderItemEntity
func (OrderItemEntity) TableName() string {
	return "order_items"
}
