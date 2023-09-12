package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept json
// @Produce json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
// @RequestHeader
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	// 调用 bindValidParam() 方法将请求中的参数绑定到 params 对象上
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// 核心逻辑
	// 1. params.Username 取得管理员信息 admininfo
	// 2. admininfo.salt + params.Password sha256 => saltPassword
	// 3. saltPassword == admininfo.password

	// 获取数据库连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 常见一个 admin 对象,用 LoginCheck 方法进行管理员登录验证
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 设置 session,并将 session 保存到会话中
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		Username:  admin.UserName,
		LoginTime: time.Now(),
	}
	// 序列化 sessInfo
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	// 传入 context 上下文，创建一个默认的会话
	sess := sessions.Default(c)
	// 将会话信息存储在 sess 中
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	fmt.Println("session", sess.Get(public.AdminSessionInfoKey))
	// 保存会话（持久化存储）
	if err := sess.Save(); err != nil {
		fmt.Printf("sess.Save() failed,err:%v\n", err)
	}

	// 成功的响应
	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
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
func (adminlogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	fmt.Println(sess.Get(public.AdminSessionInfoKey))
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
