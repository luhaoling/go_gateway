package dao

import (
	"time"

	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	// 将表单提交的数据传入数据库进行查询，判断用户是否存在
	adminInfo, err := t.Find(c, tx, &Admin{UserName: param.Username, IsDelete: 0})
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	// 传入查询到的盐以及密码，获取到含有盐值的密码字符串
	saltPassword := public.GenSaltPassword(adminInfo.Salt, param.Password)
	// 比较新生成的密码和原有密码
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码错误,请重新输入")
	}
	return adminInfo, nil

}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	// 获取 context 上下文信息后，根据传入的 search 结构体内容，到数据库中去查询相关内容，然后把结果存储到结构体中
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
