package tcp_proxy_router

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/e421083458/go_gateway/tcp_proxy_middleware"

	"github.com/e421083458/go_gateway/reverse_proxy"
	"github.com/e421083458/go_gateway/tcp_server"

	"github.com/e421083458/go_gateway/dao"
)

var tcpServerList = []*tcp_server.TcpServer{}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func TcpServerRun() {
	// 获取 Tcp 服务列表
	serviceList := dao.ServiceManagerHandler.GetTcpServiceList()
	// 遍历服务列表，初始化服务
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			// 获取服务端口
			addr := fmt.Sprintf(":%d", serviceDetail.TCPRule.Port)
			// 获取负载均衡器
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}
			// 构建路由及设置中间件
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				//tcp_proxy_middleware.IpWhiteListMiddleWare(),
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				tcp_proxy_middleware.TCPFlowLimitMiddleware(),
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)
			//
			// 构建回调handler
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(
				func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
					return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
				}, router)
			fmt.Println("Addr:", addr)
			// 构建下游地址
			baseCtx := context.WithValue(context.Background(), "service", serviceDetail)
			tcpServer := &tcp_server.TcpServer{
				Addr:    addr,
				Handler: routerHandler,
				BaseCtx: baseCtx,
			}
			// 添加初始化好的 Tcp 服务到服务列表中
			tcpServerList = append(tcpServerList, tcpServer)
			log.Printf(" [INFO] tcp_proxy_run %v\n", addr)
			// 监听端口
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Fatalf(" [INFO] tcp_proxy_run %v err:%v\n", addr, err)
			}

		}(tempItem)
	}
}

func TcpServerStop() {
	for _, tcpServer := range tcpServerList {
		tcpServer.Close()
		log.Printf(" [INFO] tcp_proxy_stop %v stopped\n", tcpServer.Addr)
	}
}
