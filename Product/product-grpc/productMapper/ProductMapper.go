package productMapper

import (
	"fmt"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"product-grpc/domain"
	"product-grpc/product"
	"product-grpc/util"
	"project1/Common/config"
	"project1/Common/contance"
)

type ProductMapper struct {
	DB *gorm.DB
}

func NewProductMapper() (*ProductMapper, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &ProductMapper{DB: db}, nil
}

func (pm *ProductMapper) ListProducts(page int, pageSize int, searchName string) ([]product.Product, int64, error) {
	var listproduct []domain.Product
	var total int64

	// 构建查询
	query := pm.DB.Model(&domain.ProductToDb{}).Where("deleted = ?", false)

	// 如果提供了搜索名称，则应用模糊搜索
	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	}

	// 计数总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&listproduct).Error; err != nil {
		return nil, 0, err
	}

	// 将 ProductToDb 转换为 Product
	var products []product.Product
	for _, pd := range listproduct {
		products = append(products, util.ConvertToProduct(pd)) // 使用转换函数
	}

	return products, total, nil
}
func (pm *ProductMapper) GetProductById(id uint32) (*product.Product, error) {
	var selectproduct domain.Product

	// 查询数据库
	if err := pm.DB.First(&selectproduct, id).Error; err != nil {
		return nil, status.Errorf(contance.USER_NOT_FOUND, "用户不存在") // 返回错误
	}

	// 将 ProductToDb 转换为 Product
	getProduct := util.ConvertToProduct(selectproduct)

	return &getProduct, nil
}

func (pm *ProductMapper) SearchProduct(column string) ([]product.Product, error) {
	// 定义一个用于存储查询结果的切片
	var getproducts []domain.Product

	// 执行查询，只选择不被标记为删除的记录
	res := pm.DB.Where("name= ?", column).Where("deleted = ?", false).Find(&getproducts)

	// 检查查询是否出错
	if res.Error != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询失败: %v", res.Error)
	}

	// 将查询结果转换为前端的 Product 列表
	var products []product.Product
	for _, pd := range getproducts {
		products = append(products, util.ConvertToProduct(pd))
	}

	return products, nil
}
func (pm *ProductMapper) SaveProduct(product *product.Product) (bool, error) {
	todb := util.ConvertToDb(product)
	res := pm.DB.Create(&todb)
	if res.Error != nil {
		return false, status.Errorf(contance.CREATE_ERROR, "新增失败")
	}
	return true, nil
}

func (pm *ProductMapper) UpdateProduct(product *product.Product) (bool, error) {
	todb := util.ConvertToDb(product)
	fmt.Println(todb.ID, todb.Name)
	res := pm.DB.Model(&domain.Product{}).Where("id = ?", todb.ID).Updates(&todb)
	if res.Error != nil {
		return false, status.Errorf(contance.UPDATE_ERROR, "更新失败")
	}
	return true, nil
}
func (pm *ProductMapper) DeleteProduct(id uint32) (bool, error) {

	res := pm.DB.Model(&domain.Product{}).Where("id =?", id).Update("deleted", true)
	if res.Error != nil {
		return false, status.Errorf(contance.DELETE_ERROR, "删除失败")
	}
	return true, nil
}
