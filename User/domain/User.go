package domain

type User struct {
	UserID      int32  `json?:"user_id,omitempty" db:"user_id" gorm:"column:user_id"` // 用户ID，前端不需要传
	Name        string `json:"name" db:"name" binding:"required"`                     // 姓名，前端必须传
	Age         int32  `json:"age" db:"age"`                                          // 年龄
	Gender      int32  `json:"gender" db:"gender"`                                    // 性别
	PhoneNumber string `json:"phone_number" db:"phone_number" binding:"required"`     // 手机号码，前端必须传
	Email       string `json:"email" db:"email"`                                      // 邮箱
	Address     string `json:"address,omitempty" db:"address"`                        // 地址，前端不需要传
	Birthday    string `json:"birthday" db:"birthday"`                                // 生日
	CreateAt    string `json:"create_at" db:"create_at"`                              // 创建时间
	UpdatedAt   string `json:"updated_at" db:"updated_at"`                            // 修改时间
	Password    string `json:"-" db:"password"`                                       // 密码，前端不需要传，生成JWT时也不包含
	Deleted     bool   `json:"-" db:"deleted"`
}

func (u User) Reset() {
	//TODO implement me
	panic("implement me")
}

func (u User) String() string {
	//TODO implement me
	panic("implement me")
}

func (u User) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}
