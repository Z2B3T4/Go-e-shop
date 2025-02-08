package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"project1/User/http/controller"
)

// UserRoutes 定义用户相关的路由
func UserRoutes(c *server.Hertz) {
	userController := controller.NewUserController()
	u := c.Group("/user")
	{

		u.GET("/getById/:userId", userController.GetById) // 直接使用指针进行方法绑定
		u.POST("/register", userController.Register)
		u.POST("/login", userController.Login)
		u.POST("/createUser", userController.CreateUser)
		u.DELETE("/deleteUser/:userId", userController.DeleteUser)
		u.POST("/logout/:userId", userController.LogOut)
		u.PUT("/update", userController.UpdateUser)
	}
}
