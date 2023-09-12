package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

// 这是一个反向代理，将请求转发给另一个服务器，并将其响应返回给客户端
// 主要实现步骤
// 创建一个 http 客户端，将上游服务器的请求头复制到该 http 客户端中
// 发起请求，获取响应（下游服务器）
// 将获取到的响应结果复制到上游服务器
func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 打印接收到的请求信息，包括请求方法、目标主机、发起请求的客户端 ip 地址和端口号
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	// 创建一个 HTTP 传输服务端
	transport := http.DefaultTransport
	// step 1,浅拷贝，然后再新增属性数据（创建一个 http.Request 的指针,然后把请求给到这个指针）
	outReq := new(http.Request)
	*outReq = *req

	// 获取客户端 ip 地址，并将其添加到请求头中的 "X-Forwarded-For" 字段中
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ",") + "," + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step 2,请求下游
	// 将请求 outReq 发送到目标服务器，并获取响应
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// step 3,把下游请求内容返回给上游
	// 遍历下游服务器的响应头,把下游服务器的响应头复制到上游服务器的响应头中
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// 给上游服务器写入状态码
	rw.WriteHeader(res.StatusCode)
	// 复制下游服务器的响应给上游服务器
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func main() {
	fmt.Println("Server on :8080")
	http.Handle("/", &Pxy{})
	http.ListenAndServe("0.0.0.0:8080", nil)
}
