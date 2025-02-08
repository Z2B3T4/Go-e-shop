package Impl

import (
	"cart-grpc/cart"
	"cart-grpc/cartMapper"
	"cart-grpc/domain"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/status"
	"os"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
)

// 提出来就可以值初始化配置一次
var cartmapper *cartMapper.CartMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context

func NewCartServiceImpl() *CartImpl {

	cartmapper, err = cartMapper.NewCartMapper()
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
	return &CartImpl{}
}

type CartImpl struct {
	cart.UnimplementedCartServiceServer
}

func (cartImpl *CartImpl) AddItem(ctx context.Context, addItemReq *cart.AddItemReq) (*cart.AddItemResp, error) {
	// 使用 mapstructure 进行对象转换
	newCartVO := domain.CartVO{}
	if err2 := mapstructure.Decode(addItemReq, &newCartVO); err2 != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "转换为CartVo失败: %v", err2)
	}
	createCart, err := cartmapper.CreateCart(&newCartVO)
	if err != nil {
		return nil, status.Errorf(contance.CREATE_ERROR, "创建购物车失败: %v", err)
	}
	if createCart == false {
		return nil, status.Errorf(contance.CREATE_ERROR, "创建购物车失败")
	}
	return &cart.AddItemResp{}, nil
}
func (cartImpl *CartImpl) GetCart(ctx context.Context, getCartReq *cart.GetCartReq) (*cart.GetCartResp, error) {
	id := getCartReq.UserId
	getCart, err := cartmapper.GetCart(int32(id))
	if err != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询购物车失败: %v", err)
	}
	return getCart, nil
}
func (cartImpl *CartImpl) EmptyCart(ctx context.Context, emptyCartReq *cart.EmptyCartReq) (*cart.EmptyCartResp, error) {
	id := emptyCartReq.UserId
	ok, err2 := cartmapper.EmptyCart(int32(id))
	if err2 != nil {
		return nil, status.Errorf(contance.DELETE_ERROR, "清空购物车失败: %v", err2)
	}
	if ok == false {
		return nil, status.Errorf(contance.DELETE_ERROR, "清空购物车失败")
	}
	return &cart.EmptyCartResp{}, nil
}
