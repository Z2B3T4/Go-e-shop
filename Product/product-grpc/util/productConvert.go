package util

import (
	"product-grpc/domain"
	"product-grpc/product"
	"strings"
)

// ConvertToDb 转换前端的 Product 为 ProductToDb
func ConvertToDb(p *product.Product) domain.Product {
	return domain.Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Picture:     p.Picture,
		Price:       p.Price,
		Deleted:     false,
		Categories:  strings.Join(p.Categories, ","), // 用逗号连接成字符串
	}
}

// ConvertToProduct 转换数据库的 ProductToDb 为前端的 Product
func ConvertToProduct(pd domain.Product) product.Product {
	return product.Product{
		Id:          pd.ID,
		Name:        pd.Name,
		Description: pd.Description,
		Picture:     pd.Picture,
		Price:       pd.Price,
		Categories:  strings.Split(pd.Categories, ","), // 用逗号分割成切片
	}
}
