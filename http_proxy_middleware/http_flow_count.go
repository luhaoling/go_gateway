package http_proxy_middleware

import (
	"fmt"
	"time"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取服务对象
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 统计项 1 全站 2 服务 3 租户
		// 获取全站统计项实例
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		// 增加全站请求数量
		totalCounter.Increase()

		// QPS:Query Per Second 每秒查询率，意思是一台服务器每秒能够响应的查询次数
		// 获取当天请求数量(全站)
		dayCount, _ := totalCounter.GetDayData(time.Now())
		fmt.Printf("totalCounter qps:%v,dayCount:%v\n", totalCounter.QPS, dayCount)
		// 获取服务统计项实例
		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		// 增加服务请求数量
		serviceCounter.Increase()

		// 获取当天请求数量
		dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		fmt.Printf("serviceCounter qps:%v,dayCount:%v\n", serviceCounter.QPS, dayServiceCount)
		c.Next()
	}
}
