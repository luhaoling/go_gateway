package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	proxy_addr = "http://127.0.0.1:2003/base"
	port       = "2002"
)

// 把访问 2002 端口的请求转发到 2003 端口 的服务器进行处理
/*
handler
请求转发，
首先是解析服务器地址（上游），然后从原有请求中获取对应参数（请求方法、请求路径和请求端口）
然后使用 http.DefaultTransport 实例化一个 http 客户端传输实例，并且使用 RoundTrip(r) 方法执行 http 请求（请求下游）得到响应
然后把请求内容返回给上游，主要是把下游的响应内容添加到响应中，返回给上游
*/
func handler(w http.ResponseWriter, r *http.Request) {
	// step 1 解析代理地址，并更改请求体的协议和主机
	proxy, err := url.Parse(proxy_addr)
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host
	r.URL.Path = proxy.Path

	// step 2 请求下游
	transport := http.DefaultTransport  // HTTP 客户端传输实例,定义了客户端与服务器之间的网络传输行为
	resp, err := transport.RoundTrip(r) // 执行一个 http 请求，返回响应和错误
	if err != nil {
		log.Print(err)
		return
	}

	// step 3 把下游请求内容返回给上游
	// 添加请求头
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()

	bufio.NewReader(resp.Body).WriteTo(w)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
