package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"payment-grpc/checkout"
	"project1/Common/CommonController"
	"project1/Common/contance"
)

func NewCheckOutController() *CheckOutController {

	addr := os.Getenv("CHECKOUT_ADDRESS")
	port := os.Getenv("CHECKOUT_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &CheckOutController{
		checkoutClient: checkout.NewCheckoutServiceClient(conn),
	}
}

type CheckOutController struct {
	CommonController.BaseController
	checkoutClient checkout.CheckoutServiceClient
}

func (checkOutController *CheckOutController) CheckOut(ctx context.Context, c *app.RequestContext) {
	req := checkout.CheckoutReq{}
	c.Bind(&req)
	resp, err := checkOutController.checkoutClient.Checkout(ctx, &req)
	if err != nil {
		checkOutController.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "创建失败", "创建失败")
		return
	}
	checkOutController.BaseController.Success(ctx, c, contance.SUCCESS, resp, "创建成功")
}
