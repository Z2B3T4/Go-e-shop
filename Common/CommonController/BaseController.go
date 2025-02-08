package CommonController

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

type BaseController struct {
}

func (b BaseController) Success(ctx context.Context, c *app.RequestContext, code int, data interface{}, msg string) {
	response := map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
	}
	c.JSON(code, response)
}

func (b BaseController) Fail(ctx context.Context, c *app.RequestContext, respcode int, code int, data interface{}, msg string) {
	response := map[string]interface{}{
		"code": code,
		"data": data,
		"msg":  msg,
	}
	c.JSON(respcode, response)

}
