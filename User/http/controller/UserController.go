package controller

import (
	"context"
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"project1/Common/CommonController"
	"project1/Common/contance"
	"project1/User/user"
	"strconv"
)

func NewUserController() *UserController {
	addr := os.Getenv("USER_ADDRESS")
	port := os.Getenv("USER_PORT")
	getenv := os.Getenv("DB_HOST")
	fmt.Println(os.Getenv("USER_ADDRESS"), os.Getenv("USER_PORT"), getenv)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("连接用户grpc失败", err)
	}

	return &UserController{
		userServiceClient: user.NewUserServiceClient(conn),
	}
}

type UserController struct {
	CommonController.BaseController
	userServiceClient user.UserServiceClient
}

// GetById retrieves a user by their ID
// @Summary Get a user by ID
// @Description Retrieve user information by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object}  nil "Success"
// @Failure 429 {object} map[string]string "Rate limit reached"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Router /user/getById/{userId} [get]
func (u *UserController) GetById(ctx context.Context, c *app.RequestContext) {
	entry, err := api.Entry("GET:/user/getById/:userId")
	if err != nil {
		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	// 获取用户ID
	userId := c.Param("userId")
	id, err2 := strconv.Atoi(userId)
	if err2 != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	//fmt.Println(id)
	// 准备请求参数
	req := user.UserInfoReq{UserID: int32(id)}
	resp, err3 := u.userServiceClient.GetUserInfo(ctx, &req)
	fmt.Printf("%+v", resp)
	if err3 != nil {
		fmt.Println("获取用户失败", err)
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.USER_NOT_FOUND, "用户未找到", "用户未找到")
		return
	}
	// 返回 JSON 响应
	u.BaseController.Success(ctx, c, contance.SUCCESS, resp, "成功")

}

func (u *UserController) Register(ctx context.Context, c *app.RequestContext) {
	entry, limit := api.Entry("POST:/user/register")
	if limit != nil {
		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	newUser := user.RegisterReq{}
	err := c.Bind(&newUser)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}

	// 调用注册方法并传递上下文
	register, err := u.userServiceClient.Register(ctx, &newUser)
	if register == nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.USER_ALREADY_EXIST, "用户已存在", "用户已存在")
		return
	}
	if err != nil {
		fmt.Println("调用注册方法失败", err)
		return // 需要适当地处理错误，比如返回一个错误响应
	}
	// 创建用户信息的响应
	userInfo := map[string]string{
		"token":  register.Token,
		"userId": string(register.UserId),
	}
	fmt.Printf("%+v", userInfo)
	u.BaseController.Success(ctx, c, contance.SUCCESS, userInfo, "注册成功")
	return
}
func (u *UserController) Login(ctx context.Context, c *app.RequestContext) {
	entry, limit := api.Entry("POST:/user/login")
	if limit != nil {

		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	newUser := user.LoginReq{}
	err := c.Bind(&newUser)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	// 调用注册方法并传递上下文
	login, err := u.userServiceClient.Login(ctx, &newUser)
	if err != nil {
		fmt.Println("调用注册方法失败", err)
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "注册失败", "注册失败")
		return // 需要适当地处理错误，比如返回一个错误响应
	}
	// 创建用户信息的响应
	userInfo := map[string]string{
		"token":  login.Token,
		"userId": strconv.Itoa(int(login.UserId)),
	}
	fmt.Printf("%+v", userInfo)
	u.BaseController.Success(ctx, c, contance.SUCCESS, userInfo, "注册成功")
	return
}
func (u *UserController) GetUserInfo(ctx context.Context, c *app.RequestContext) {
	strid := c.Param("userId")
	id, err := strconv.Atoi(strid)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
	}
	userinfo, err := u.userServiceClient.GetUserInfo(ctx, &user.UserInfoReq{UserID: int32(id)})
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.SELECT_ERROR, "获取失败", "获取失败")
	}
	u.BaseController.Success(ctx, c, contance.SUCCESS, userinfo, "获取成功")
}
func (u *UserController) CreateUser(ctx context.Context, c *app.RequestContext) {
	entry, limit := api.Entry("POST:/user/createUser")
	if limit != nil {
		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	newUser := user.CreateUserReq{}
	err := c.Bind(&newUser)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	createUser, err := u.userServiceClient.CreateUser(ctx, &newUser)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.USER_CREATE_ERROR, "创建失败", "创建失败")
		return
	}
	u.BaseController.Success(ctx, c, contance.SUCCESS, createUser, "创建成功")
}
func (u *UserController) DeleteUser(ctx context.Context, c *app.RequestContext) {
	entry, limit := api.Entry("DELETE:/user/deleteUser/:userId")
	if limit != nil {
		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	strId := c.Param("userId")
	intid, err := strconv.Atoi(strId)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	deleteUser, err := u.userServiceClient.DeleteUser(ctx, &user.DeleteUserReq{UserID: int32(intid)})
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.DELETE_ERROR, "删除失败", "删除失败")
		return
	}
	u.BaseController.Success(ctx, c, contance.SUCCESS, deleteUser, "删除成功")
}
func (u *UserController) UpdateUser(ctx context.Context, c *app.RequestContext) {
	entry, limit := api.Entry("PUT:/user/update")
	if limit != nil {
		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	newUser := user.UpdateUserReq{}
	err := c.Bind(&newUser)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	createUser, err := u.userServiceClient.UpdateUser(ctx, &newUser)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.UPDATE_ERROR, "更新失败", "更新失败")
		return
	}
	u.BaseController.Success(ctx, c, contance.SUCCESS, createUser, "更新成功")
}
func (u *UserController) LogOut(ctx context.Context, c *app.RequestContext) {
	entry, limit := api.Entry("POST:/user/logout/:userId")
	if limit != nil {
		// 被限流的处理逻辑
		u.BaseController.Fail(ctx, c, contance.LIMIT_ERROR, contance.LIMIT_ERROR, "请求限流", "请求限流")
		return
	}
	defer entry.Exit()
	strId := c.Param("userId")
	intid, err := strconv.Atoi(strId)
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.PARAM_ERROR, "参数错误", "参数错误")
		return
	}
	out, err := u.userServiceClient.LogOut(ctx, &user.LogOutReq{UserId: int32(intid)})
	if err != nil {
		u.BaseController.Fail(ctx, c, contance.PARAM_ERROR, contance.LOGOUT_ERROR, "退出失败", "退出失败")
	}
	u.BaseController.Success(ctx, c, contance.SUCCESS, out, "退出成功")
}
