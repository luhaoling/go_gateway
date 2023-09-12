package controller

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type OAuthController struct{}

func OAuthRegister(group *gin.RouterGroup) {
	oauth := &OAuthController{}
	group.POST("/tokens", oauth.Tokens)
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]

func (oauth *OAuthController) Tokens(c *gin.Context) {
	params := &dto.TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 划分租户账密
	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		middleware.ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	//fmt.Println(splits)
	// 对租户账密进行 Base64 解码(Basic Auth 是一种 http 身份验证方法，它使用 Base64 编码)
	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//fmt.Println("appSecret", string(appSecret))

	// 取出 app_id secret
	// 生成 app_list
	// 匹配 app_id
	// 基于 jwt 生成 token
	// 生成 output
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		middleware.ResponseError(c, 2003, errors.New("用户或密码名格式错误"))
		return
	}

	appList := dao.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,                                                               // 发行人
				ExpiresAt: time.Now().Add(public.JwtExpired * time.Second).In(lib.TimeLocation).Unix(), // JWT 的过期时间
			}
			// 生成 token
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 2004, err)
				return
			}
			output := &dto.TokensOutput{
				ExpiresIn:   public.JwtExpired, // 过期时间
				TokenType:   "Bearer",          // token 类型（Bearer Token 用于进行身份验证和授权）
				AccessToken: token,             // token
				Scope:       "read_write",      // 权限范围
			}

			middleware.ResponseSuccess(c, output)
			return
		}
	}
	middleware.ResponseError(c, 2005, errors.New("未匹配正确的 APP 信息"))
}

// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
// @Security ApiKeyAuth
func (adminlogin *OAuthController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	fmt.Println(sess.Get(public.AdminSessionInfoKey))
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
