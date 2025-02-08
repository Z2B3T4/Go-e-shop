package userMapper

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"project1/Common/config"
	"project1/Common/contance"
	"project1/User/domain"
)

type UserMapper struct {
	DB *gorm.DB
}

func NewUserMapper() (*UserMapper, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := cfg.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &UserMapper{DB: db}, nil
}

func (um UserMapper) Save(user *domain.User) int32 {
	//err := um.DB.AutoMigrate(&domain.User{})
	//if err != nil {
	//	fmt.Println("创建数据库失败")
	//	return -1
	//}
	preUser, _ := um.GetByEmail(user.Email)
	if preUser != nil {
		return -1
	}
	um.DB.Create(&user)
	getUser, _ := um.GetByEmail(user.Email)
	fmt.Printf("%#v", getUser.UserID)
	return getUser.UserID
}
func (um *UserMapper) GetById(id int) (*domain.User, error) {
	var user domain.User // 假设 User 是你在 auth_domain 包中定义的用户模型

	if err := um.DB.Where("deleted = ?", 0).First(&user, id).Error; err != nil {
		return nil, err
	}

	claim := &domain.User{
		UserID:      user.UserID,
		Name:        user.Name,
		Age:         user.Age,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Address:     user.Address,
		Birthday:    user.Birthday,
		CreateAt:    user.CreateAt,
		UpdatedAt:   user.UpdatedAt,
		Password:    user.Password,
	}

	return claim, nil
}

func (um *UserMapper) GetByEmail(email string) (*domain.User, error) {
	var user domain.User // 假设 User 是你在 auth_domain 包中定义的用户模型
	if err := um.DB.Where("email = ?", email).Where("deleted=?", 0).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.Unimplemented, "根据邮箱查询异常")
	}

	claim := &domain.User{
		UserID:      user.UserID,
		Name:        user.Name,
		Age:         user.Age,
		Gender:      user.Gender,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Address:     user.Address,
		Birthday:    user.Birthday,
		CreateAt:    user.CreateAt,
		UpdatedAt:   user.UpdatedAt,
		Password:    user.Password,
	}

	return claim, nil
}

func (um *UserMapper) DeleteUser(id int32) (bool, error) {

	// 更新用户记录
	result := um.DB.Model(&domain.User{}).Where("user_id = ?", id).Where("deleted=?", 0).Update("Deleted", true)
	if result.Error != nil {
		fmt.Println("删除失败", result.Error)
		return false, status.Errorf(contance.SELECT_ERROR, "删除失败(用户可能已经被删除)")
	} else if result.RowsAffected == 0 {
		fmt.Println(" 没有找到要删除的记录")
		return false, status.Errorf(contance.USER_NOT_FOUND, "用户未找到")
	}
	return true, nil
}

func (um *UserMapper) UpdateUser(user domain.User) (bool, error) {
	result := um.DB.Model(&user).Where("user_id = ?", user.UserID).Where("deleted=?", 0).Updates(user)
	if result.Error != nil {
		fmt.Println("更新失败", result.Error)
		return false, status.Errorf(contance.SELECT_ERROR, "删除失败")
	} else if result.RowsAffected == 0 {
		fmt.Println(" 没有找到要更新的记录或者没有更新任何记录")
		return false, status.Errorf(contance.USER_NOT_FOUND, "用户未找到")
	}
	return true, nil
}
