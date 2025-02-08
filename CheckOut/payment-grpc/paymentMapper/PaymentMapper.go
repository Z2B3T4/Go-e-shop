package paymentMapper

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"math/rand"
	"order-grpc/order"
	"os"
	"payment-grpc/checkout"
	"payment-grpc/domain"
	"payment-grpc/domain/VO"
	"project1/Common/config"
	"project1/Common/contance"
	"strconv"
	"time"
)

type PaymentMapper struct {
	DB          *gorm.DB
	orderClient order.OrderServiceClient
}

func NewPaymentMapper() (*PaymentMapper, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	addr := os.Getenv("ORDER_ADDRESS")
	port := os.Getenv("ORDER_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}
	return &PaymentMapper{DB: db, orderClient: order.NewOrderServiceClient(conn)}, nil
}

func (pm *PaymentMapper) Checkout(req *checkout.CheckoutReq, orderResp *order.ListOrderResp) (*checkout.CheckoutResp, error) {
	var paymentInfo domain.Payment
	var paymenttemp VO.CreditCardInfo
	if err := mapstructure.Decode(req.CreditCard, &paymenttemp); err != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "failed to decode RegisterReq to User: %v", err)
	}

	// 检查用户信息
	flag := false
	for _, item := range orderResp.Orders {
		//same, err := AreEqual(item.Address, req.Address)
		//if err != nil {
		//	return nil, err
		//}
		if item.Email == req.Email {
			if Validate(req.CreditCard) == false {
				return nil, status.Errorf(contance.PARAM_ERROR, "信用卡信息不合法")
			}
			// 使用 mapstructure 进行对象转换
			if err := mapstructure.Decode(paymenttemp, &paymentInfo); err != nil {
				return nil, status.Errorf(contance.CONVERT_ERROR, "failed to decode RegisterReq to User: %v", err)
			}
			paymentInfo.UserID = item.UserId
			paymentInfo.OrderID = item.OrderId
			timestamp := time.Now().Unix()
			randomSuffix := rand.Intn(100) // 生成 0 到 99 之间的随机数
			paymentInfo.TransactionID = fmt.Sprintf("ORD-%d-%02d", timestamp, randomSuffix)
			flag = true
		}
	}
	if flag == false {
		return nil, status.Errorf(contance.PARAM_ERROR, "邮箱错误")
	}
	result := pm.DB.Model(&domain.Payment{}).Create(&paymentInfo)
	if result.Error != nil {
		return nil, status.Errorf(contance.CREATE_ERROR, "创建支付信息失败"+result.Error.Error())
	}
	resp := &checkout.CheckoutResp{
		OrderId:       paymentInfo.OrderID,
		TransactionId: paymentInfo.TransactionID,
	}
	return resp, nil

}

// AreEqual 比较两个相同类型的结构体是否相等
func AreEqual(a *order.Address, b *checkout.Address) (bool, error) {
	if a.StreetAddress != b.StreetAddress {
		return false, nil
	}
	if a.City != b.City {
		return false, nil
	}
	if a.State != b.State {
		return false, nil
	}
	if a.Country != b.Country {
		return false, nil
	}

	if strconv.Itoa(int(a.ZipCode)) != b.ZipCode {
		return false, nil
	}

	// 如果所有字段相等，返回 true
	return true, nil
}

// Validate 检查信用卡信息的有效性
func Validate(cc *checkout.CreditCardInfo) bool {
	// 检查信用卡号码的格式
	if len(cc.CreditCardNumber) < 13 || len(cc.CreditCardNumber) > 19 {
		return false
	}
	for _, ch := range cc.CreditCardNumber {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	// 检查 CVV 的值范围
	if cc.CreditCardCvv < 100 || cc.CreditCardCvv > 9999 { // 3 到 4 位数
		return false
	}

	// 检查到期日期
	now := time.Now()
	expirationDate := time.Date(int(cc.CreditCardExpirationYear), time.Month(cc.CreditCardExpirationMonth), 1, 0, 0, 0, 0, time.UTC)
	if expirationDate.Before(now) {
		return false
	}

	return true
}

func (pm *PaymentMapper) GetPaymentItem(userId uint32, OrderId uint32) (checkout.GetPaymentItemResp, error) {
	getItem := &domain.Payment{}
	result := pm.DB.Model(&domain.Payment{}).Where("order_id = ? and user_id = ?", OrderId, userId).First(&getItem)
	if result.Error != nil {
		return checkout.GetPaymentItemResp{}, status.Errorf(contance.SELECT_ERROR, "查询支付信息失败"+result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return checkout.GetPaymentItemResp{}, status.Errorf(contance.SELECT_ERROR, "支付信息不存在")
	}
	resp := checkout.GetPaymentItemResp{}
	// 使用 mapstructure 进行对象转换
	if err := mapstructure.Decode(getItem, &resp); err != nil {
		return checkout.GetPaymentItemResp{}, status.Errorf(contance.CONVERT_ERROR, "failed to decode RegisterReq to User: %v", err)
	}
	return resp, nil
}
