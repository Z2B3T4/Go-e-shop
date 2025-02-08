package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"project1/Auth/auth"
	jwtUtil "project1/Auth/jwtUtil"
	"project1/Common/config"
	"project1/Common/contance"
	"project1/User/domain"
	"strconv"
	"time"
)

// 提出来就可以值初始化配置一次
// var cartmapper *cartMapper.CartMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context

func NewAuthServiceImpl() *TokenServiceImpl {

	if err != nil {
		fmt.Println(err, "状态码：", contance.SERVER_ERROR)

	}
	// 初始化redis
	db := os.Getenv("REDIS_DB")
	dbnum, err2 := strconv.Atoi(db)
	if err2 != nil {
		fmt.Println(err2)
	}
	redisClient = config.NewRedisClient(dbnum) // 指定第一个分区
	return &TokenServiceImpl{}
}

type TokenServiceImpl struct {
	auth.UnimplementedAuthServiceServer
}
type MyClaims struct {
	UserId               int32
	jwt.RegisteredClaims // 这个RegisteredClaims是一个接口，里面有jwt的
}

func (t *TokenServiceImpl) DeliverTokenByRPC(ctx context.Context, deliverTokenReq *auth.DeliverTokenReq) (*auth.DeliveryResp, error) {
	hs := jwtUtil.HS{
		Key:        os.Getenv("JWT_KEY"),
		SignMethod: jwtUtil.HS256,
	}
	claims := MyClaims{
		UserId: deliverTokenReq.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	encode, err := hs.Encode(claims)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "生成token失败: %v", err)
	}
	resp := auth.DeliveryResp{Token: encode}
	return &resp, nil
}
func (t *TokenServiceImpl) VerifyTokenByRPC(ctx context.Context, verifyTokenReq *auth.VerifyTokenReq) (*auth.VerifyResp, error) {
	// 配置文件中读取存入redis的第几分区
	stringdb := os.Getenv("REDIS_DB")
	intdb, err := strconv.Atoi(stringdb)
	if err != nil {
		return &auth.VerifyResp{Res: false}, status.Errorf(contance.PARAM_ERROR, "读取redis分区失败")
	}
	// 初始化redis
	var redisClient = config.NewRedisClient(intdb)
	// 判断token是否存在

	exists := redisClient.Exists(contance.TOKEN_REDIS_PREFIX + verifyTokenReq.Token)
	if exists.Val() <= 0 {
		return &auth.VerifyResp{Res: false}, status.Errorf(contance.TOKEN_VERIFY_ERROR, "用户不存在,token校验失败")
	}
	// 解析token
	var getUser domain.MyClaim
	//getUser := domain.MyClaim{}
	// 读取配置文件
	//godotenv.Load("./config.env")
	hs := jwtUtil.HS{Key: os.Getenv("JWT_KEY"), SignMethod: jwtUtil.HS256}
	err2 := hs.Decode(verifyTokenReq.Token, &getUser)
	// 解析失败
	if err2 != nil {
		return &auth.VerifyResp{Res: false}, status.Errorf(contance.TOKEN_VERIFY_ERROR, "token解析失败: "+err2.Error())
	}
	return &auth.VerifyResp{Res: true}, nil

}
