package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"http/controller"
)

func CartRoutes(h *server.Hertz) {
	orderController := controller.NewOrderController()
	u := h.Group("/order")
	{
		u.POST("/createOrder", orderController.CreateOrder)
		u.GET("/listOrder/:userId", orderController.GetOrder)
		u.PUT("/markOrder", orderController.MarkOrder)
	}
}
