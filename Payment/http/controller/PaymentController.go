package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"paymentGrpc/payment"
	"project1/Common/CommonController"
	"project1/Common/contance"
)

func NewPaymentController() *PaymentController {

	addr := os.Getenv("PAYMENT_ADDRESS")
	port := os.Getenv("PAYMENT_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &PaymentController{
		paymentClient: payment.NewPaymentServiceClient(conn),
	}
}

type PaymentController struct {
	CommonController.BaseController
	paymentClient payment.PaymentServiceClient
}

func (pc *PaymentController) Charge(ctx context.Context, c *app.RequestContext) {
	var req payment.ChargeReq
	if err := c.Bind(&req); err != nil {
		pc.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "参数错误", "参数错误")
		return
	}
	res, err := pc.paymentClient.Charge(ctx, &req)
	if err != nil {
		pc.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.CREATE_ERROR, "支付失败"+err.Error(), "支付失败")
		return
	}
	pc.BaseController.Success(ctx, c, contance.SUCCESS, res, "支付成功")

}
