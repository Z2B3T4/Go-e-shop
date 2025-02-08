package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"http/controller"
)

func CartRoutes(h *server.Hertz) {
	checkOutController := controller.NewCheckOutController()
	u := h.Group("/checkout")
	{
		u.POST("/checkout", checkOutController.CheckOut)
	}
}
