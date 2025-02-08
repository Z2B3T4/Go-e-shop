package jwtV1

import "github.com/golang-jwt/jwt/v4"

type JwtValidator interface {
	Encode(claims jwt.Claims) (string, error) // 这个Claims就是一个接口，里面是生成jwt时所带的数据
	Decode(sign string, claims jwt.Claims) error
}
