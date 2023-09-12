package dao

import (
	"time"

	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdateAt    time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreateAt    time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除 0:否 1:是"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := t.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}
	// 查询 HttpRule 的内容
	httpRule := &HttpRule{ServiceID: search.ID}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 查询 TcpRule 的内容
	tcpRule := &TcpRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 查询 GrpcRule 的内容
	grpcRule := &GrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 查询 AccessControl 的内容
	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 查询 LoadBalance 的内容
	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

func (t *ServiceInfo) GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
	list := []dto.DashServiceStatItemOutput{}
	query := tx.SetCtx(public.GetGinTraceContext(c))
	// 使用 where("is_delete=0") 方法添加了一个查询条件
	// 使用 select("load_type,count(*) as value") 指定要查询的字段
	// 使用 Group("load_type") 指定了分组字段
	if err := query.Table(t.TableName()).Where("is_delete=0").Select("load_type,count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (t *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize

	// 获取 gin 上下文，并且将其设置在 gorm 查询中
	query := tx.SetCtx(public.GetGinTraceContext(c))
	// 设置查询只返回 is_delete=0 的数据
	query = query.Table(t.TableName()).Where("is_delete=0")
	// 判断是否需要进行模糊查询
	if param.Info != "" {
		// 对 service_name 和 service_desc 这两列进行模糊查询
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	// param.PageSize 指定每页的数量,Offset 是偏移量（即从哪一行开始返回数据）, id desc 是按照 id 这一列进行降序排序
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, nil
	}
	// 计算符合查询条件的总数
	//total = int64(len(list))
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	// 进行查询,如果查询到内容，则 err 为 nil;
	// 如果查询不到内容，则 err 不为 nil
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
