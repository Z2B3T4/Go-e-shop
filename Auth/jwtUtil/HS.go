package jwtV1

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
)

type HS struct {
	Key        string
	SignMethod HSSignMethod
}
type HSSignMethod string

const (
	HS256 HSSignMethod = "HS256"
	HS384 HSSignMethod = "HS384"
	HS512 HSSignMethod = "HS512"
)

func (hs *HS) getSignMethod() *jwt.SigningMethodHMAC {
	switch hs.SignMethod {
	case HS256:
		return jwt.SigningMethodHS256
	case HS384:
		return jwt.SigningMethodHS384
	case HS512:
		return jwt.SigningMethodHS512
	}
	return jwt.SigningMethodHS256
}

// Encode 生成令牌字符串。
// 该方法接收一个jwt.Claims接口类型的参数claims，用于指定令牌的声明。
// 返回值是一个字符串类型的令牌和一个错误类型，用于表示编码过程中遇到的错误（如果有）。
func (hs *HS) Encode(claims jwt.Claims) (string, error) {
	// 创建一个新的JWT令牌，并使用指定的签名方法和声明。
	token := jwt.NewWithClaims(hs.getSignMethod(), claims)

	// 使用HS实例的密钥对令牌进行签名并生成签名后的令牌字符串。
	// 如果生成过程中出现错误，则记录错误并返回空字符串和错误。
	sign, err := token.SignedString([]byte(hs.Key))
	if err != nil {
		log.Fatal(err)
		return sign, err
	}

	// 如果没有错误，返回签名后的令牌字符串和nil作为错误值。
	return sign, nil
}

// Decode 验证JWT令牌并解析其声明。
// 该方法使用HS算法验证签名，并将声明填充到claims参数中。
// 参数:
//
//	sign: 需要解码的JWT令牌字符串。
//	claims: 一个实现jwt.Claims接口的对象，用于存储解析后的令牌声明。
//
// 返回值:
//
//	如果解析过程中发生错误，则返回该错误。
func (hs *HS) Decode(sign string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(sign, claims, func(token *jwt.Token) (interface{}, error) {
		// 使用HS实例的Key作为密钥来验证令牌。
		return []byte(hs.Key), nil
	})
	return err
}
