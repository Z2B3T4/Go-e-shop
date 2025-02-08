package cartMapper

import (
	"cart-grpc/cart"
	"cart-grpc/domain"
	"cart-grpc/util"
	"fmt"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"project1/Common/config"
	"project1/Common/contance"
)

type CartMapper struct {
	DB *gorm.DB
}

func NewCartMapper() (*CartMapper, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("购物车配置文件加载失败: %w", err)
	}

	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("购物车链接mysql数据库失败: %w", err)
	}

	return &CartMapper{DB: db}, nil
}

func (cm *CartMapper) CreateCart(cartVo *domain.CartVO) (bool, error) {
	newcart := util.ConvertToDb(cartVo)
	result := cm.DB.Model(&domain.Cart{}).Create(&newcart)
	if result.Error != nil {
		return false, status.Errorf(contance.SELECT_ERROR, "添加购物车失败")
	}
	return true, nil
}

// GetCart 根据用户 ID 获取购物车信息
func (cm *CartMapper) GetCart(userID int32) (*cart.GetCartResp, error) {
	var carts []domain.Cart

	// 查询用户的购物车记录
	result := cm.DB.Model(&domain.Cart{}).Where("user_id = ? AND deleted = ?", userID, false).Find(&carts)
	if result.Error != nil {
		return nil, result.Error
	}
	// 将购物车记录转换为 GetCartResp 格式
	cartResp := &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: uint32(userID),
			Items:  make([]*cart.CartItem, 0),
		},
	}

	// 遍历所有购物车记录，填充 CartItem
	for _, getcart := range carts {
		cartResp.Cart.Items = append(cartResp.Cart.Items, &cart.CartItem{
			ProductId: uint32(getcart.ProductID),
			Quantity:  getcart.Quantity,
		})
	}

	return cartResp, nil
}

func (cm *CartMapper) EmptyCart(userId int32) (bool, error) {
	result := cm.DB.Model(&domain.Cart{}).Where("user_id =?", userId).Where("deleted=?", 0).Update("deleted", true)
	if result.Error != nil {
		return false, status.Errorf(contance.SERVER_ERROR, "清空购物车查询数据库失败")
	}
	if result.RowsAffected == 0 {
		return false, status.Errorf(contance.SERVER_ERROR, "清空购物车失败（没有查询到）")
	}
	return true, nil
}
