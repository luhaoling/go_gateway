package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs1.Run()
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	// syscall.SIGINT ctrl+c
	// syscall.SIGTERM kill()函数发送
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type RealServer struct {
	Addr string
}

func (r *RealServer) Run() {
	log.Println("Starting httpserver at +" + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error/", r.ErrorHandler)
	mux.HandleFunc("test_http_string/test_http_string/aaa/", r.TimeoutHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

// Handler 和 server 的区别和联系
// 区别：
// Handler 是一个实现了 ServeHTTP 方法的接口，主要用来处理请求和返回响应
// Server 是一个创建和配置服务器的结构体。主要包含服务器的属性和启动服务器、关闭服务器的方法
// 联系：
// Server 是创建和配置服务器的实例，需要一个 Handler 来处理请求，以提供完整的 http 服务器功能

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	upath := fmt.Sprintf("http://%s%s", r.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded_For=%v,X-Real_IP=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-ip"))
	headers := fmt.Sprintf("headers=%v\n", req.Header)
	// X-Forwarded-For :记录经过的每一层代理（包含整个代理链路上的所有 ip），当前代理地址由下一个代理或服务器添加
	// x-real-for:存储的是真实的客户端 ip 地址
	// x-forwarded-for 和 x-real-for 存储的ip，都是由代理服务器添加
	// RemoteAddr 是 go 标准库的一个字段。上诉两个均是 http 请求头部的字段

	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, headers)
	fmt.Println("path", upath)
	for i, v := range req.Header {
		fmt.Printf("%v:%v\n", i, v)
	}
}

func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *RealServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(7 * time.Second)
	updath := "time headler"
	w.WriteHeader(200)
	io.WriteString(w, updath)
}
