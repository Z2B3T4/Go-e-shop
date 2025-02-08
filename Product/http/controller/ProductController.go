package controller

import (
	"context"
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"product-grpc/product"
	"project1/Common/CommonController"
	"project1/Common/contance"
	"strconv"
)

func NewProductController() *ProductController {
	addr := os.Getenv("PRODUCT_ADDRESS")
	port := os.Getenv("PRODUCT_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &ProductController{
		productClient: product.NewProductCatalogServiceClient(conn),
	}
}

type ProductController struct {
	CommonController.BaseController
	productClient product.ProductCatalogServiceClient
}

func (p *ProductController) GetList(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("GET:/product/getList")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	column := c.Query("column")
	var intpage int
	var err error
	if intpage, err = strconv.Atoi(page); err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	var intpagesize int
	if intpagesize, err = strconv.Atoi(pageSize); err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	req := product.ListProductsReq{
		Page:         int32(intpage),
		PageSize:     int64(intpagesize),
		CategoryName: column,
	}
	res, err := p.productClient.ListProducts(ctx, &req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "查询失败", "查询失败")
		return
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "查询成功")

}

func (p *ProductController) GetById(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("GET:/product/getById/:productId")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	productId := c.Param("productId")
	id, err := strconv.Atoi(productId)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	req := product.GetProductReq{Id: uint32(id)}
	res, err := p.productClient.GetProduct(ctx, &req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "查询失败", "查询失败")
		return
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "查询成功")
}

func (p *ProductController) SearchByName(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("GET:/product/getByName")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	query := c.Query("query")
	req := product.SearchProductsReq{Query: query}
	res, err := p.productClient.SearchProducts(ctx, &req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "查询失败", "查询失败")
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "查询成功")
}

func (p *ProductController) CreateProduct(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("POST:/product/create")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	var newproduct product.Product
	err := c.Bind(&newproduct)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	req := product.CreateProductReq{
		Product: &newproduct,
	}
	res, err := p.productClient.CreateProduct(ctx, &req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "创建失败", "创建失败")
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "创建成功")
}

func (p *ProductController) UpdateProduct(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("PUT:/product/update")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
	}
	defer entry.Exit()
	var newproduct product.Product
	err := c.Bind(&newproduct)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	req := product.UpdateProductReq{
		Product: &newproduct,
	}
	res, err := p.productClient.UpdateProduct(ctx, &req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.UPDATE_ERROR, "更新失败", "更新失败")
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "成功")
}

func (p *ProductController) DeleteProduct(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("DELETE:/product/delete/:productId")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
	}
	defer entry.Exit()
	productId := c.Param("productId")
	id, err := strconv.Atoi(productId)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	req := product.DeleteProductReq{Id: uint32(id)}
	res, err := p.productClient.DeleteProduct(ctx, &req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.DELETE_ERROR, "删除失败", "删除失败")
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "成功")
}
