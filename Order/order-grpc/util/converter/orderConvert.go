package converter

import (
	"order-grpc/domain/Entity"
	"order-grpc/domain/VO"
	"order-grpc/order"
	"strconv"
)

func ConvertPlaceOrderReqToAddressEntity(req VO.PlaceOrderReq) Entity.AddressEntity {
	return Entity.AddressEntity{
		UserID:        int32(req.UserID),
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		Country:       req.Address.Country,
		ZipCode:       req.Address.ZipCode,
		Deleted:       false,
	}
}

func ConvertPlaceOrderReqToOrderEntity(req VO.PlaceOrderReq, addressID int32) Entity.OrderEntity {
	return Entity.OrderEntity{
		UserID:       int32(req.UserID),
		AddressID:    addressID,
		UserCurrency: req.UserCurrency,
		Email:        req.Email,
		Deleted:      false,
	}
}

func ConvertPlaceOrderReqToOrderItemEntities(req VO.PlaceOrderReq, orderId int32) []Entity.OrderItemEntity {
	var orderItemEntities []Entity.OrderItemEntity
	for _, item := range req.OrderItems {
		orderItemEntities = append(orderItemEntities, Entity.OrderItemEntity{
			OrderID:   orderId,
			ProductID: int32(item.Item.ProductID),
			Quantity:  item.Item.Quantity,
			Cost:      item.Cost,
			Deleted:   false,
		})
	}
	return orderItemEntities
}

func ConvertAddressEntityToProto(address Entity.AddressEntity) order.Address {
	return order.Address{
		StreetAddress: address.StreetAddress,
		City:          address.City,
		State:         address.State,
		Country:       address.Country,
		ZipCode:       address.ZipCode,
	}
}

func ConvertOrderEntityToProto(orderEntity Entity.OrderEntity, address order.Address) order.Order {
	return order.Order{
		OrderId:      strconv.Itoa(int(orderEntity.ID)),
		UserId:       uint32(orderEntity.UserID),
		UserCurrency: orderEntity.UserCurrency,
		Email:        orderEntity.Email,
		Address:      &address,
	}
}

func ConvertOrderItemEntityToProto(orderItem Entity.OrderItemEntity) order.OrderItem {
	return order.OrderItem{
		Item: &order.CartItem{
			ProductId: uint32(orderItem.ProductID),
			Quantity:  orderItem.Quantity,
		},
		Cost: orderItem.Cost,
	}
}

/*
	func ConvertToListOrderResp(address Entity.AddressEntity, orderEntity Entity.OrderEntity, orderItems []Entity.OrderItemEntity) order.ListOrderResp {
		var orderItemProtos []*order.OrderItem // 使用指针切片

		// 将 OrderItemEntity 转换为 OrderItem proto 消息并添加到切片中
		for _, item := range orderItems {
			orderItemProto := ConvertOrderItemEntityToProto(item)
			orderItemProtos = append(orderItemProtos, &orderItemProto) // 将指针添加到切片
		}

		// 创建 orderProto，包含装载了对应 address 和 order items 的信息
		orderProto := ConvertOrderEntityToProto(orderEntity, ConvertAddressEntityToProto(address))
		orderProto.OrderItems = orderItemProtos // 赋值指针切片

		// 创建指向 orderProto 的指针并返回
		return order.ListOrderResp{
			Orders: []*order.Order{&orderProto}, // 使用切片指向 order.Order 指针
		}
	}
*/
func ConvertToListOrderResp(addressMap map[int32]Entity.AddressEntity, orderEntities []Entity.OrderEntity, orderItemsMap map[int32][]Entity.OrderItemEntity) order.ListOrderResp {
	var listOrders []*order.Order // 用于存储多个 order.Order 的指针切片

	// 遍历所有订单实体
	for _, orderEntity := range orderEntities {
		// 从 map 中获取与订单 ID 相关的所有订单项
		orderItems := orderItemsMap[orderEntity.ID]

		// 将 OrderItemEntity 转换为 OrderItem proto 消息并添加到切片中
		var orderItemsProtos []*order.OrderItem
		for _, item := range orderItems {
			orderItemProto := ConvertOrderItemEntityToProto(item)
			orderItemsProtos = append(orderItemsProtos, &orderItemProto) // 将指针添加到切片
		}

		// 根据地址 ID 从 addressMap 获取地址信息
		address, addressExists := addressMap[orderEntity.AddressID]
		if !addressExists {
			// 如果我们找不到地址，处理错误或返回空地址结构
			// 这里对地址缺失的处理可以根据业务需求调整
			address = Entity.AddressEntity{} // 或者可以直接跳过这个订单
		}

		// 创建 orderProto，包含装载了对应 address 和 order items 的信息
		orderProto := ConvertOrderEntityToProto(orderEntity, ConvertAddressEntityToProto(address))
		orderProto.OrderItems = orderItemsProtos // 赋值订单项指针切片

		// 将 orderProto 指针添加到 listOrders
		listOrders = append(listOrders, &orderProto)
	}

	// 创建并返回 listOrderResp，包含多个订单
	return order.ListOrderResp{
		Orders: listOrders, // 返回多个 order.Order 指针
	}
}
