package http_proxy_middleware

import (
	"fmt"
	"strings"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// jwt_auth_token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//fmt.Println("serviceDetail", serviceDetail)
		//decode jwt token
		// app_id 与 app_list 取得 appInfo
		// 使用空字符串替换掉 token 前面的Bearer
		token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		fmt.Println("token", token)
		appMatched := false
		if token != "" {
			claims, err := public.JwtDecode(token)
			if err != nil {
				middleware.ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			//fmt.Println("claims.Issuer", claims.Issuer)
			// claims.Issuer 存储的是 AppID
			appList := dao.AppManagerHandler.GetAppList()
			for _, appInfo := range appList {
				if appInfo.AppID == claims.Issuer {
					// 验证通过后，将租户信息存进 context 中
					c.Set("app", appInfo)
					//fmt.Println(appInfo)
					appMatched = true
					break
				}
			}
			if serviceDetail.AccessControl.OpenAuth == 1 && !appMatched {
				middleware.ResponseError(c, 2003, errors.New("not match valid app"))
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
