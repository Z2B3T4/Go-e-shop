package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"http/controller"
)

func CartRoutes(h *server.Hertz) {
	cartController := controller.NewCartController()
	u := h.Group("/cart")
	{
		u.POST("/addCart", cartController.AddItem)
		u.GET("/getCart/:userId", cartController.GetItem)
		u.PUT("/emptyCart/:userId", cartController.EmptyCart)
	}
}
