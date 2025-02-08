package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"http/controller"
)

func PaymentRoutes(h *server.Hertz) {
	paymentController := controller.NewPaymentController()
	u := h.Group("/payment")
	{
		u.POST("/create", paymentController.Charge)
	}
}
