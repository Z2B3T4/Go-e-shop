package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"product-grpc/product"
	"product-grpc/productMapper"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
)

// 提出来就可以值初始化配置一次
var productmapper *productMapper.ProductMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context

func NewProductService() *ProductImpl {

	productmapper, err = productMapper.NewProductMapper()
	if err != nil {
		fmt.Println(err, "状态码：", contance.SERVER_ERROR)

	}
	// 初始化redis
	db := os.Getenv("REDIS_DB")
	dbnum, err2 := strconv.Atoi(db)
	if err2 != nil {
		fmt.Println(err2)
	}
	redisClient = config.NewRedisClient(dbnum) // 指定第一个分区
	return &ProductImpl{}
}

type ProductImpl struct {
	product.UnimplementedProductCatalogServiceServer
}

func (productImpl *ProductImpl) ListProducts(ctx context.Context, listProductsReq *product.ListProductsReq) (*product.ListProductsResp, error) {
	products, _, err := productmapper.ListProducts(int(listProductsReq.Page), int(listProductsReq.PageSize), listProductsReq.CategoryName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}
	// 将 products 转换为 []*Product
	var productPointers []*product.Product
	for _, p := range products {
		productPointers = append(productPointers, &p) // 取每个元素的指针
	}
	res := &product.ListProductsResp{Products: productPointers}
	return res, nil
}
func (productImpl *ProductImpl) GetProduct(ctx context.Context, getProductReq *product.GetProductReq) (*product.GetProductResp, error) {
	getProduct, err2 := productmapper.GetProductById(getProductReq.Id)
	if err2 != nil {
		return nil, status.Errorf(contance.USER_NOT_FOUND, "用户不存在: %v", err2)
	}
	resp := product.GetProductResp{Product: getProduct}
	return &resp, nil
}
func (productImpl *ProductImpl) SearchProducts(ctx context.Context, searchProductsReq *product.SearchProductsReq) (*product.SearchProductsResp, error) {
	searchProduct, err := productmapper.SearchProduct(searchProductsReq.Query)
	if err != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询失败: %v", err)
	}
	var productPointers []*product.Product
	for _, p := range searchProduct {
		productPointers = append(productPointers, &p) // 取每个元素的指针
	}
	res := &product.SearchProductsResp{Results: productPointers}
	return res, nil
}
func (productImpl *ProductImpl) CreateProduct(ctx context.Context, createProduct *product.CreateProductReq) (*product.CreateProductResp, error) {
	reqProduct := createProduct.Product
	ok, err := productmapper.SaveProduct(reqProduct)
	if err != nil {
		return nil, status.Errorf(contance.CREATE_ERROR, "创建失败: %v", err)
	}
	if ok == false {
		return nil, status.Errorf(contance.CREATE_ERROR, "创建失败")
	}
	res := &product.CreateProductResp{Success: ok}
	return res, nil
}
func (productImpl *ProductImpl) UpdateProduct(ctx context.Context, updateProduct *product.UpdateProductReq) (*product.UpdateProductResp, error) {
	reqProduct := updateProduct.Product
	fmt.Println(reqProduct.Name, " awd")
	ok, err := productmapper.UpdateProduct(reqProduct)
	if err != nil {
		return nil, status.Errorf(contance.UPDATE_ERROR, "修改失败: %v", err)
	}
	if ok == false {
		return nil, status.Errorf(contance.UPDATE_ERROR, "修改失败")
	}
	res := &product.UpdateProductResp{Success: ok}
	return res, nil
}
func (productImpl *ProductImpl) DeleteProduct(ctx context.Context, deleteProduct *product.DeleteProductReq) (*product.DeleteProductResp, error) {
	ok, err := productmapper.DeleteProduct(deleteProduct.Id)
	if err != nil {
		return nil, status.Errorf(contance.DELETE_ERROR, "删除失败: %v", err)
	}
	if ok == false {
		return nil, status.Errorf(contance.DELETE_ERROR, "删除失败")
	}
	res := &product.DeleteProductResp{Success: ok}
	return res, nil
}
