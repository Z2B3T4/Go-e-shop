package routes

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"http/controller"
)

func ProductRoutes(c *server.Hertz) {
	productController := controller.NewProductController()
	u := c.Group("/product")
	{
		u.GET("/getList", productController.GetList)
		u.GET("/getById/:productId", productController.GetById)
		u.GET("/getByName", productController.SearchByName)
		u.POST("/create", productController.CreateProduct)
		u.DELETE("/delete/:productId", productController.DeleteProduct)
		u.PUT("/update", productController.UpdateProduct)
	}
}
