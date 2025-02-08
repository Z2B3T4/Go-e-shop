package domain

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

type MyClaim struct {
	UserId int32
	//MyUser  User
	Myclaim jwt.RegisteredClaims
}

func (c *MyClaim) Valid() error {
	// 调用 Embedded struct 的 Valid() 方法进行基本的验证
	if err := c.Myclaim.Valid(); err != nil {
		return err // 返回标准的 Token 验证错误
	}

	// 自定义逻辑：可以在这里添加其他验证
	// 例如，检查用户 ID 是否有效
	if c.UserId <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	// 如果一切正常返回 nil
	return nil
}
