package Impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure" // 对象转换
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	jwtUtil "project1/Auth/jwtUtil"
	"project1/Common/config"
	"project1/Common/contance"
	"project1/User/domain"
	"project1/User/user"
	"project1/User/userMapper"
	"strconv"
	"time"
)

// 提出来就可以值初始化配置一次
var usermapper *userMapper.UserMapper
var err error
var redisClient *redis.Client
var CenterContext *context.Context

func NewUserService() *UserServiceImpl {
	// 获取当前工作目录
	//cwd, err := os.Getwd()
	//if err != nil {
	//	log.Fatalf("Error getting current working directory: %v", err)
	//}
	//fmt.Println(cwd)
	//// 构建上下层目录的绝对路径
	//envFilePath := filepath.Join(cwd, "User", "userhttp.env") // 确保路径正确
	//fmt.Println(envFilePath)
	//if err := godotenv.Load(envFilePath); err != nil {
	//	log.Fatalf("Error loading env file: %v", err)
	//}

	usermapper, err = userMapper.NewUserMapper()
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
	return &UserServiceImpl{}
}

type UserServiceImpl struct {
	user.UnimplementedUserServiceServer
}

func (userController *UserServiceImpl) Register(ctx context.Context, registerReq *user.RegisterReq) (*user.RegisterResp, error) {
	newUser := domain.User{}
	fmt.Printf("%+v", registerReq)

	// 使用 mapstructure 进行对象转换
	if err := mapstructure.Decode(registerReq, &newUser); err != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "failed to decode RegisterReq to User: %v", err)
	}
	id := usermapper.Save(&newUser)
	if id == -1 {
		return nil, status.Errorf(contance.USER_CREATE_ERROR, "用户已经存在")
	}
	myclaim := domain.MyClaim{
		//MyUser: newUser,
		UserId: id,
		Myclaim: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	key := os.Getenv("JWT_KEY")
	hs := jwtUtil.HS{Key: key, SignMethod: jwtUtil.HS256}
	sign, err := hs.Encode(&myclaim)
	if err != nil {
		fmt.Println(err, "生成token失败")
		return nil, status.Errorf(contance.TOKEN_GENERATE_ERROR, "生成token失败: %v", err)
	}

	fmt.Println(sign)

	// 将 token 设置到 Metadata 中
	md := metadata.New(map[string]string{"token": sign})
	grpc.SetHeader(ctx, md)
	// token 生成完成

	// 存入 redis
	strid := strconv.Itoa(int(id))
	redisClient.Set(contance.TOKEN_REDIS_PREFIX+sign, strid, contance.REDIS_TOKEN_EXPIRE_TIME)
	redisClient.Set(contance.TOKEN_REDIS_PREFIX+strid, sign, contance.REDIS_TOKEN_EXPIRE_TIME)

	registerResp := &user.RegisterResp{}
	registerResp.UserId = id
	registerResp.Token = sign

	return registerResp, nil
}

func (userController *UserServiceImpl) Login(ctx context.Context, loginReq *user.LoginReq) (*user.LoginResp, error) {
	getUser, err2 := usermapper.GetByEmail(loginReq.Email)
	if getUser == nil {
		return nil, status.Errorf(contance.USER_NOT_FOUND, "用户不存在")
	}
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询用户失败")
	}

	// 判断密码
	if getUser.Password != loginReq.Password {
		return nil, status.Errorf(contance.LOGIN_EMAIL_PASSWORD_ERROR, "邮箱或者密码错误")
	}
	// 生成jwt的Claim
	myclaim := &domain.MyClaim{
		//MyUser: *getUser,
		UserId: getUser.UserID,
		Myclaim: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// 指定jwt的秘钥和加密方式
	hs := jwtUtil.HS{Key: os.Getenv("JWT_KEY"), SignMethod: jwtUtil.HS256}
	token, err2 := hs.Encode(myclaim)
	if err2 != nil {
		fmt.Println(err2, "登录生成token失败")
		return nil, status.Errorf(contance.TOKEN_GENERATE_ERROR, "登录生成token失败: %v", err2)
	}
	// 存入redis
	redisClient.Set(token, string(getUser.UserID), contance.REDIS_TOKEN_EXPIRE_TIME)
	redisClient.Set(contance.TOKEN_REDIS_PREFIX+string(getUser.UserID), token, contance.REDIS_TOKEN_EXPIRE_TIME)

	// 将用户信息和token写入上下文
	context.WithValue(ctx, "userId", getUser.UserID)
	context.WithValue(ctx, "userName", getUser.Name)
	context.WithValue(ctx, "token", token)

	res := &user.LoginResp{UserId: getUser.UserID, Token: token}
	return res, nil
}
func (userController *UserServiceImpl) GetUserInfo(ctx context.Context, userInfoReq *user.UserInfoReq) (*user.UserInfoResp, error) {
	getUser, err2 := usermapper.GetById(int(userInfoReq.UserID))
	if getUser == nil {
		return nil, status.Errorf(contance.USER_NOT_FOUND, "用户不存在")
	}
	if err2 != nil {
		return nil, status.Errorf(contance.SELECT_ERROR, "查询用户失败")
	}
	resp := &user.UserInfoResp{}
	// 使用 mapstructure 进行对象转换
	if err := mapstructure.Decode(getUser, &resp); err != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "failed to decode RegisterReq to User: %v", err)
	}
	return resp, nil
}
func (userController *UserServiceImpl) CreateUser(ctx context.Context, createUserReq *user.CreateUserReq) (*user.CreateUserResp, error) {
	newUser := &domain.User{}
	if err := mapstructure.Decode(createUserReq, &newUser); err != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "转换失败: %v", err)
	}
	id := usermapper.Save(newUser)
	if id == -1 {
		return &user.CreateUserResp{Ok: false}, status.Errorf(contance.USER_CREATE_ERROR, "用户创建失败,用户已存在")
	}
	if id >= 0 {
		return &user.CreateUserResp{Ok: true}, nil
	}
	return &user.CreateUserResp{Ok: false}, status.Errorf(contance.USER_CREATE_ERROR, "用户创建失败")
}
func (userController *UserServiceImpl) DeleteUser(ctx context.Context, deleteUserReq *user.DeleteUserReq) (*user.DeleteUserResp, error) {
	ok, err := usermapper.DeleteUser(deleteUserReq.UserID)
	if err != nil {
		return &user.DeleteUserResp{Ok: false}, err
	}
	if ok == true {
		return &user.DeleteUserResp{Ok: true}, nil
	}
	return &user.DeleteUserResp{Ok: false}, status.Errorf(contance.USER_NOT_FOUND, "用户删除失败，用户不存在")
}
func (userController *UserServiceImpl) LogOut(ctx context.Context, logOutReq *user.LogOutReq) (*user.LogOutResp, error) {
	del := redisClient.Del(contance.TOKEN_REDIS_PREFIX + string(logOutReq.UserId))
	if del.Err() != nil {
		fmt.Println("用户缓存删除失败")
		return &user.LogOutResp{Ok: false}, status.Errorf(contance.REDIS_DELETE_ERROR, "用户缓存删除失败")
	}
	if del.Val() == 0 {
		return &user.LogOutResp{Ok: false}, status.Errorf(contance.USER_NOT_FOUND, "用户不存在")

	}
	return &user.LogOutResp{Ok: true}, nil
}
func (userController *UserServiceImpl) UpdateUser(ctx context.Context, updateUserReq *user.UpdateUserReq) (*user.UpdateUserResp, error) {
	newUser := &domain.User{}
	if err := mapstructure.Decode(updateUserReq, &newUser); err != nil {
		return nil, status.Errorf(contance.CONVERT_ERROR, "转换失败: %v", err)
	}
	ok, err2 := usermapper.UpdateUser(*newUser)
	if err2 != nil {
		return nil, status.Errorf(contance.UPDATE_ERROR, "更新失败: %v", err2)
	}
	if ok {
		return &user.UpdateUserResp{Ok: true}, nil
	}
	return &user.UpdateUserResp{Ok: false}, status.Errorf(contance.USER_CREATE_ERROR, "用户创建失败")
}
