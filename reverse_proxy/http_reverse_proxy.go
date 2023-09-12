package reverse_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/e421083458/go_gateway/middleware"

	"github.com/e421083458/go_gateway/reverse_proxy/load_balance"
	"github.com/gin-gonic/gin"
)

func NewLoadBalanceReverseProxy(c *gin.Context, lb load_balance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	// 请求协调者
	direction := func(req *http.Request) {
		// 获取目标地址
		nextAddr, err := lb.Get(req.URL.String())
		// todo 优化点3
		if err != nil || nextAddr == "" {
			panic("get next addr fail")
		}
		// 解析这个目标地址
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		// 目标地址和请求地址的结合
		req.URL.Path = singleJoiningSSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User_Agent"]; !ok {
			req.Header.Set("User-Agent", "user_agent")
		}
	}

	// 更改内容
	modifyFunc := func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		//todo 优化点2
		//var payload []byte
		//var readErr error
		//
		//if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		//	gr, err := gzip.NewReader(resp.Body)
		//	if err != nil {
		//		return err
		//	}
		//	payload, readErr = ioutil.ReadAll(gr)
		//	resp.Header.Del("Content-Encoding")
		//} else {
		//	payload, readErr = ioutil.ReadAll(resp.Body)
		//}
		//if readErr != nil {
		//	return readErr
		//}
		//
		//c.Set("status_code", resp.StatusCode)
		//c.Set("payload", payload)
		//resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
		//resp.ContentLength = int64(len(payload))
		//resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
		return nil
	}

	// 错误回调：关闭 real_server 时测试，错误回调
	// 范围：transport.RoundTrip 发生的错误、以及 ModifyResponse 发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		middleware.ResponseError(c, 999, err)
	}
	return &httputil.ReverseProxy{Director: direction, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}

func singleJoiningSSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
