package router

import (
	"context"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler *http.Server
)

func HttpServerRun() {
	// 根据配置文件设置 gin 框架运行模式
	gin.SetMode(lib.GetStringConf("base.base.debug_mode"))
	// 初始化路由，返回一个路由实例
	r := InitRouter()
	// 根据配置创建一个 http.Server 对象
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("base.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("base.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("base.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.http.max_header_bytes")),
	}
	// 开启一个新的 goroutine,执行 http 服务器的启动逻辑
	go func() {
		// 打印服务器启动的信息和监听的地址
		log.Printf(" [INFO] HttpServerRun:%s\n", lib.GetStringConf("base.http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			log.Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", lib.GetStringConf("base.http.addr"), err)
		}
	}()
}

func HttpServerStop() {
	// 创建一个带有超时时间的上下文，设置超时时间为 10s
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 优雅关机
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpServerStop stopped\n")
}
