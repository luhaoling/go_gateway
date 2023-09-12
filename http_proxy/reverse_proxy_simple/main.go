package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	// 127.0.0.1:2002/xxx
	// 127.0.0.1:2003/base/xxx
	rs1 := "http://127.0.0.1:2003/base/error/"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}

	// 创建一个反向代理的 http 处理器（创建一个正向到单个主机的反向代理处理器）
	proxy := httputil.NewSingleHostReverseProxy(url1)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
