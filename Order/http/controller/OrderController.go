package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"order-grpc/order"
	"os"
	"project1/Common/CommonController"
	"project1/Common/contance"
	"strconv"
)

func NewOrderController() *OrderController {

	addr := os.Getenv("ORDER_ADDRESS")
	port := os.Getenv("ORDER_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &OrderController{
		orderClient: order.NewOrderServiceClient(conn),
	}
}

type OrderController struct {
	CommonController.BaseController
	orderClient order.OrderServiceClient
}

func (orderController *OrderController) CreateOrder(ctx context.Context, c *app.RequestContext) {
	createreq := order.PlaceOrderReq{}
	err := c.Bind(&createreq)
	if err != nil {
		orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "创建失败", "创建失败")
		return
	}
	res, err := orderController.orderClient.PlaceOrder(ctx, &createreq)
	if err != nil {
		orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "创建失败", "创建失败")
		return
	}
	orderController.BaseController.Success(ctx, c, contance.SUCCESS, res, "创建成功")
}

func (orderController *OrderController) GetOrder(ctx context.Context, c *app.RequestContext) {
	strid := c.Param("userId")
	intid, err := strconv.Atoi(strid)
	if err != nil {
		orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	req := &order.ListOrderReq{UserId: uint32(intid)}
	res, err := orderController.orderClient.ListOrder(ctx, req)
	if err != nil {
		orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "查询失败", "查询失败")
		return
	}
	orderController.BaseController.Success(ctx, c, contance.SUCCESS, res, "查询成功")
}

func (orderController *OrderController) MarkOrder(ctx context.Context, c *app.RequestContext) {
	req := order.MarkOrderPaidReq{}
	//param := c.Query("user_id")
	//intuserid, err2 := strconv.Atoi(param)
	//if err2 != nil {
	//	orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误2", "参数错误")
	//	return
	//}
	//req.UserId = uint32(intuserid)
	//strorderId := c.Query("order_id")
	//req.OrderId = strorderId
	err2 := c.Bind(&req)
	fmt.Printf("%+v", &req)
	if err2 != nil {
		orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	res, err := orderController.orderClient.MarkOrderPaid(ctx, &req)
	if err != nil {
		orderController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "查询失败"+err.Error(), "查询失败")
		return
	}
	orderController.BaseController.Success(ctx, c, contance.SUCCESS, res, "查询成功")
}
