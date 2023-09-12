package http_proxy_middleware

import (
	"errors"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/reverse_proxy"
	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取对应服务信息
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("server not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 获取对应的负载均衡器
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		trans, err := dao.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		//middleware.ResponseSuccess(c,"ok")
		//return

		// 创建 reverseproxy
		// 使用 reverseproxy.ServerHTTP(c.Request,c.Response)
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}

}
