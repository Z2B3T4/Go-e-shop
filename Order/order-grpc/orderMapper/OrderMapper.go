package orderMapper

import (
	"fmt"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"order-grpc/domain/Entity"
	"order-grpc/domain/VO"
	"order-grpc/order"
	"order-grpc/util/converter"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
)

type OrderMapper struct {
	DB *gorm.DB
}

func NewOrderMapper() (*OrderMapper, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &OrderMapper{DB: db}, nil
}

func (om *OrderMapper) CreateOrder(req *VO.PlaceOrderReq) (int32, error) {
	addressEntity := converter.ConvertPlaceOrderReqToAddressEntity(*req)
	result := om.DB.Model(&Entity.AddressEntity{}).Create(&addressEntity)
	if result.Error != nil {
		return 0, status.Errorf(contance.CREATE_ERROR, "创建地址失败")
	}
	addressId := addressEntity.ID
	/*
		id, ok := result.Get("ID")
		if ok == false {
			return 0, status.Errorf(contance.CREATE_ERROR, "获取地址ID失败")
		}
	*/
	orderEntity := converter.ConvertPlaceOrderReqToOrderEntity(*req, addressId)
	result = om.DB.Model(&Entity.OrderEntity{}).Create(&orderEntity)
	if result.Error != nil {
		return 0, status.Errorf(contance.CREATE_ERROR, "创建订单失败"+result.Error.Error())
	}
	orderId := orderEntity.ID
	orderItemEntities := converter.ConvertPlaceOrderReqToOrderItemEntities(*req, orderId)
	result = om.DB.Model(&Entity.OrderItemEntity{}).Create(&orderItemEntities)
	if result.Error != nil {
		return 0, status.Errorf(contance.CREATE_ERROR, "创建订单条目失败"+result.Error.Error())
	}
	return orderId, nil
}
func (om *OrderMapper) ListOrder(userId int32) (*order.ListOrderResp, error) {
	// 查找用户所有有效订单信息
	var orderEntities []Entity.OrderEntity
	result := om.DB.Model(&Entity.OrderEntity{}).Where("user_id = ? ", userId).Where("deleted = ?", 0).Find(&orderEntities)
	if result.Error != nil {
		return nil, status.Errorf(contance.SERVER_ERROR, "用户订单不存在")
	}

	// 使用 map 来存储地址 ID 到地址的映射
	addressMap := make(map[int32]Entity.AddressEntity)

	// 查找所有订单的地址 ID
	addressIDs := make([]int32, 0)
	for _, orderEntity := range orderEntities {
		addressIDs = append(addressIDs, orderEntity.AddressID)
	}

	// 根据地址 ID 查询地址信息
	var addresses []Entity.AddressEntity
	result = om.DB.Model(&Entity.AddressEntity{}).Where("id IN ?", addressIDs).Where("deleted = ?", 0).Find(&addresses)
	if result.Error != nil {
		return nil, status.Errorf(contance.SERVER_ERROR, "用户地址未找到")
	}

	// 将地址存储到 map 中
	for _, addr := range addresses {
		addressMap[addr.ID] = addr
	}

	// 使用 map 来存储订单 ID 和对应的订单项
	orderItemsMap := make(map[int32][]Entity.OrderItemEntity)

	// 查找每个订单的订单项信息
	for _, orderEntity := range orderEntities {
		var orderItems []Entity.OrderItemEntity
		result = om.DB.Model(&Entity.OrderItemEntity{}).Where("order_id = ?", orderEntity.ID).Where("deleted = ?", 0).Find(&orderItems)
		if result.Error != nil {
			return nil, status.Errorf(contance.SERVER_ERROR, "订单 ID %d 下的订单条目不存在", orderEntity.ID)
		}
		// 将订单项存储在 map 中
		orderItemsMap[orderEntity.ID] = orderItems
	}

	// 组装返回信息
	resp := converter.ConvertToListOrderResp(addressMap, orderEntities, orderItemsMap)

	return &resp, nil
}

func (om *OrderMapper) MarkOrder(userId int32, orderId int32) (bool, error) {
	// 开始一个新的事务
	tx := om.DB.Begin()

	// 标记订单为删除
	result := tx.Model(&Entity.OrderEntity{}).
		Where("user_id = ? ", userId).
		Where("deleted = ?", 0).
		Update("deleted", 1)

	flag := false
	listOrder, err := om.ListOrder(userId)
	if err != nil {
		tx.Rollback()
		return false, status.Errorf(contance.SERVER_ERROR, "查询订单失败")
	}
	for _, item := range listOrder.Orders {
		intid, err := strconv.Atoi(item.OrderId)
		if err != nil {
			return false, status.Errorf(contance.SERVER_ERROR, "订单ID转换失败")
		}
		if int32(intid) == int32(orderId) {
			flag = true
		}
	}
	if flag == false {
		tx.Rollback()
		return false, status.Errorf(contance.SERVER_ERROR, "订单不在改用户下，不存在")
	}
	// 标记订单项为删除
	result2 := tx.Model(&Entity.OrderItemEntity{}).
		Where("order_id = ?", orderId).
		Where("deleted = ?", 0).
		Update("deleted", 1)

	// 检查操作结果
	if result.Error != nil || result2.Error != nil || result.RowsAffected == 0 || result2.RowsAffected == 0 {
		tx.Rollback()
		return false, status.Errorf(contance.SERVER_ERROR, "标记订单失败")
	}

	// 提交事务
	tx.Commit()
	return true, nil
}
