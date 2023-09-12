package tcp_proxy_middleware

import (
	"fmt"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/public"
)

func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 统计项 1 全站 2 服务 3 租户
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		totalCounter.Increase()

		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		fmt.Println("TCP流量统计,totalCount", totalCounter.TotalCount, "QPS", totalCounter.QPS)
		fmt.Println("TCP流量统计,serviceCounter", serviceCounter.TotalCount, "QPS", serviceCounter.QPS)
		c.Next()
	}
}
