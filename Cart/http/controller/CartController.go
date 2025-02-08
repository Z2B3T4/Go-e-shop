package controller

import (
	"cart-grpc/cart"
	"context"
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"http/domain"
	"os"
	"project1/Common/CommonController"
	"project1/Common/contance"
	"strconv"
)

func NewCartController() *CartController {

	addr := os.Getenv("CART_ADDRESS")
	port := os.Getenv("CART_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &CartController{
		cartClient: cart.NewCartServiceClient(conn),
	}
}

type CartController struct {
	CommonController.BaseController
	cartClient cart.CartServiceClient
}

func (p *CartController) AddItem(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("POST:/cart/AddCart")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	cartVO := domain.CartVO{}
	err := c.Bind(&cartVO)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	if cartVO.UserID <= 0 || cartVO.Item.ProductID <= 0 || cartVO.Item.Quantity <= 0 {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	req := &cart.AddItemReq{
		UserId: cartVO.UserID,
		Item: &cart.CartItem{
			ProductId: cartVO.Item.ProductID,
			Quantity:  cartVO.Item.Quantity,
		},
	}
	item, err := p.cartClient.AddItem(ctx, req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "添加购物车失败", "添加购物车失败")
		return
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, item, "添加购物车成功")

}

func (p *CartController) GetItem(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("GET:/cart/GetCart")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	strUserId := c.Param("userId")
	intUserId, err := strconv.Atoi(strUserId)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	req := &cart.GetCartReq{
		UserId: uint32(intUserId),
	}
	res, err := p.cartClient.GetCart(ctx, req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "查询购物车失败", "查询购物车失败")
		return
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "查询购物车成功")
}

func (p *CartController) EmptyCart(ctx context.Context, c *app.RequestContext) {
	entry, err3 := api.Entry("PUT:/cart/emptyCart/:userId")
	if err3 != nil {
		// 被限流的处理逻辑
		p.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	strUserId := c.Param("userId")
	intUserId, err := strconv.Atoi(strUserId)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	req := &cart.EmptyCartReq{
		UserId: uint32(intUserId),
	}
	res, err := p.cartClient.EmptyCart(ctx, req)
	if err != nil {
		p.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.DELETE_ERROR, "清空购物车失败", "清空购物车失败")
		return
	}
	p.BaseController.Success(ctx, c, contance.SUCCESS, res, "清空购物车成功")
}
