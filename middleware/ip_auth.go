package middleware

import (
	"errors"
	"fmt"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/gin-gonic/gin"
)

// ip 地址认证的 gin 中间件。
// 它会检查请求的客户端 ip 地址是否在配置文件中指定的允许访问的 ip 地址列表，如果不在列表中，则返回一个错误响应，并且中止请求处理
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _, host := range lib.GetStringSliceConf("base.http.allow_ip") {
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			ResponseError(c, InternalErrorCode, errors.New(fmt.Sprintf("%v, not in iplist", c.ClientIP())))
			c.Abort()
			return
		}
		c.Next()
	}
}
