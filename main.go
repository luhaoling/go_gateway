package main

import (
	"flag"
	"github.com/e421083458/go_gateway/grpc_proxy_router"
	"os/signal"
	"syscall"

	"github.com/e421083458/go_gateway/tcp_proxy_router"

	"github.com/e421083458/go_gateway/http_proxy_router"

	"os"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/e421083458/go_gateway/router"
)

// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /

// endpoint dashboard 后台管理 server 代理服务器
// config ./conf/prod/ 对应配置文件夹

var (
	endpoint = flag.String("endpoint", "server", "input endpoint dashboard or server")
	config   = flag.String("config", "./conf/dev", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 根据命令行输入的 dashboard 还是 server 判断是进入哪个模式
	if *endpoint == "dashboard" {
		lib.InitModule(*config)
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		lib.InitModule(*config)
		defer lib.Destroy()
		// 服务信息初始化（存储在 map 和 slice 中）
		dao.ServiceManagerHandler.LoadOnce()
		// 租户信息初始化（存储在 map 和 slice 中）
		dao.AppManagerHandler.LoadOnce()

		go func() {
			http_proxy_router.HttpServerRun()
		}()

		go func() {
			http_proxy_router.HttpsServerRun()
		}()

		go func() {
			tcp_proxy_router.TcpServerRun()
		}()
		go func() {
			grpc_proxy_router.GrpcServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		grpc_proxy_router.GrpcServerStop()
		tcp_proxy_router.TcpServerStop()
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()

	}

}
