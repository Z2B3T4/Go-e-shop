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
	checkout "payment-grpc/checkout"
	"payment-grpc/paymentMapper"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
)

// 提出来就可以值初始化配置一次
var paymentmapper *paymentMapper.PaymentMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context
var orderClient order.OrderServiceClient

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
	addr := os.Getenv("ORDER_ADDRESS")
	port := os.Getenv("ORDER_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	orderClient = order.NewOrderServiceClient(conn)
	return &PaymentImpl{}
}

type PaymentImpl struct {
	checkout.UnimplementedCheckoutServiceServer
}

func (paymentImpl *PaymentImpl) Checkout(ctx context.Context, checkoutReq *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	listOrder, err2 := orderClient.ListOrder(ctx, &order.ListOrderReq{UserId: checkoutReq.UserId})
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询订单失败service: %v", err2)
	}
	if len(listOrder.Orders) == 0 {
		return nil, status.Errorf(contance.SELECT_ERROR, "订单为空")
	}
	resp, err2 := paymentmapper.Checkout(checkoutReq, listOrder)
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "支付失败service: %v", err2)
	}
	return resp, nil

}
func (paymentImpl *PaymentImpl) GetPaymentItem(ctx context.Context, getPaymentItemReq *checkout.GetPaymentItemReq) (*checkout.GetPaymentItemResp, error) {
	intorderid, err2 := strconv.Atoi(getPaymentItemReq.OrderId)
	if err2 != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "标记订单失败service: %v", err)
	}
	item, err2 := paymentmapper.GetPaymentItem(getPaymentItemReq.UserId, uint32(intorderid))
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询订单失败service: %v", err2)
	}
	return &item, nil
}
