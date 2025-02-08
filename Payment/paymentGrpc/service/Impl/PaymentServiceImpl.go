package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"order-grpc/order"
	"os"
	"payment-grpc/checkout"
	"paymentGrpc/payment"
	"paymentGrpc/paymentMapper"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
)

// 提出来就可以值初始化配置一次
var paymentmapper *paymentMapper.PaymentMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context

func NewPaymentService() *PaymentImpl {

	paymentmapper, err = paymentMapper.NewPaymentMapper()
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

	addr := os.Getenv("CHECKOUT_ADDRESS")
	port := os.Getenv("CHECKOUT_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	addr2 := os.Getenv("ORDER_ADDRESS")
	port2 := os.Getenv("ORDER_PORT")
	conn2, err2 := grpc.Dial(fmt.Sprintf("%s:%s", addr2, port2), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err2 != nil {
		fmt.Println("连接用户grpc失败", err2)
	}
	return &PaymentImpl{
		checkoutClient: checkout.NewCheckoutServiceClient(conn),
		orderClient:    order.NewOrderServiceClient(conn2),
	}
}

type PaymentImpl struct {
	payment.UnimplementedPaymentServiceServer
	checkoutClient checkout.CheckoutServiceClient
	orderClient    order.OrderServiceClient
}

func (paymentImpl *PaymentImpl) Charge(ctx context.Context, chargeReq *payment.ChargeReq) (*payment.ChargeResp, error) {
	resp, err2 := paymentImpl.checkoutClient.GetPaymentItem(ctx, &checkout.GetPaymentItemReq{UserId: chargeReq.UserId, OrderId: chargeReq.OrderId})
	if err2 != nil {
		return nil, status.Errorf(contance.CREATE_ERROR, "支付失败"+err2.Error())
	}
	create, err3 := paymentmapper.Create(chargeReq, resp)
	if err3 != nil {
		return nil, status.Errorf(contance.CREATE_ERROR, "创建支付信息失败"+err3.Error())
	}
	_, err2 = paymentImpl.orderClient.MarkOrderPaid(ctx, &order.MarkOrderPaidReq{UserId: chargeReq.UserId, OrderId: chargeReq.OrderId})
	if err2 != nil {
		return nil, status.Errorf(contance.CREATE_ERROR, "标记完成订单失败"+err2.Error())
	}
	result := payment.ChargeResp{TransactionId: create.TransactionID}
	return &result, nil
}
