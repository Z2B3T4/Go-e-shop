package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/status"
	"order-grpc/domain/VO"
	"order-grpc/order"
	"order-grpc/orderMapper"
	"os"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
)

// 提出来就可以值初始化配置一次
var ordermapper *orderMapper.OrderMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context

func NewOrderService() *OrderImpl {

	ordermapper, err = orderMapper.NewOrderMapper()
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
	return &OrderImpl{}
}

type OrderImpl struct {
	order.UnimplementedOrderServiceServer
}

func (orderImpl *OrderImpl) PlaceOrder(ctx context.Context, placeOrderReq *order.PlaceOrderReq) (*order.PlaceOrderResp, error) {
	// 使用 mapstructure 进行对象转换
	newplaceOrderreq := VO.PlaceOrderReq{}
	if err := mapstructure.Decode(placeOrderReq, &newplaceOrderreq); err != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "failed to decode RegisterReq to User: %v", err)
	}
	orderId, err2 := ordermapper.CreateOrder(&newplaceOrderreq)
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "创建订单失败service: %v", err2)
	}
	strId := strconv.Itoa(int(orderId))
	resp := order.PlaceOrderResp{Order: &order.OrderResult{OrderId: strId}}
	return &resp, nil
}
func (orderImpl *OrderImpl) ListOrder(ctx context.Context, listOrderReq *order.ListOrderReq) (*order.ListOrderResp, error) {
	listOrder, err2 := ordermapper.ListOrder(int32(listOrderReq.UserId))
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询订单失败service: %v", err2)
	}
	return listOrder, nil
}
func (orderImpl *OrderImpl) MarkOrderPaid(ctx context.Context, markOrderPaidReq *order.MarkOrderPaidReq) (*order.MarkOrderPaidResp, error) {
	intOrderID, err2 := strconv.Atoi(markOrderPaidReq.OrderId)
	if err2 != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "标记订单失败service: %v", err)
	}
	ok, err2 := ordermapper.MarkOrder(int32(markOrderPaidReq.UserId), int32(intOrderID))
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "标记订单失败service: %v", err2)
	}
	if ok == false {
		return nil, status.Errorf(contance.SELECT_ERROR, "标记订单失败service: %v", err2)
	}
	return &order.MarkOrderPaidResp{}, nil
}
