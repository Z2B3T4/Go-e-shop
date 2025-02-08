package domain

import "time"

type ProductToDb struct {
	ID          uint32    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Picture     string    `json:"picture"`
	Price       float32   `gorm:"not null" json:"price"`
	Deleted     bool      `gorm:"default:false" json:"deleted"`
	Categories  string    `gorm:"not null" json:"categories"` // 使用逗号分隔的字符串
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 返回绑定的数据库表名
func (ProductToDb) TableName() string {
	return "products" // 指定表名为 products
}
