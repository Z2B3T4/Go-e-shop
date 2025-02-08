package contance

import "time"

type ExpireTime struct{}

var _ = ExpireTime{} // 确保 ExpireTime 结构体被使用

const (
	REDIS_TOKEN_EXPIRE_TIME = 60 * time.Minute // 1天
)
