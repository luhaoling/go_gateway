package http_proxy_middleware

import (
	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 查看是否有符合的前缀和域名
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			// 终止，终止后面所有的该请求下的函数
			c.Abort()
			return
		}
		//fmt.Println("matched service", public.Obj2Json(service))
		// 查询到后存储对应的 service
		c.Set("service", service)
		c.Next()
	}

}
