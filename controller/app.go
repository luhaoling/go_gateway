package controller

import (
	"time"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type APPController struct{}

// APPControllerRegister admin 路由注册
func APPRegister(router *gin.RouterGroup) {
	admin := APPController{}
	router.GET("app_list", admin.APPList)
	router.GET("app_detail", admin.APPDetail)
	router.GET("app_stat", admin.APPStatistics)
	router.GET("app_delete", admin.APPDelete)
	router.POST("app_add", admin.APPAdd)
	router.POST("app_update", admin.APPUpdate)
}

// APPList godoc
// @Summary 租户列表
// @Description 租户列表
// @Tags 租户管理
// @ID /app/app_list
// @Accept json
// @Product json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
// @Router /app/app_list [get]
func (admin *APPController) APPList(c *gin.Context) {
	// 获取请求参数
	params := &dto.APPListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 查询租户列表
	info := &dao.App{}
	list, total, err := info.APPList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 租户列表字段格式化
	outputList := []dto.APPListItemOutput{}
	for _, item := range list {
		appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + item.AppID)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		outputList = append(outputList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  appCounter.TotalCount,
			RealQps:  appCounter.QPS,
		})
	}
	output := dto.APPListOutput{
		List:  outputList,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// APPDetail godoc
// @Summary 租户详情
// @Description 租户详情
// @Tags 租户管理
// @ID /app/app_detail
// @Accept json
// @Product json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=dto.APP} "success"
// @Router /app/app_detail [get]
func (admin *APPController) APPDetail(c *gin.Context) {
	// 获取请求参数
	params := &dto.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, detail)
	return
}

// APPAdd godoc
// @Summary 租户添加
// @Description 租户添加
// @Tags 租户管理
// @ID /app/app_add
// @Accept json
// @Product json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_add [post]
func (admin *APPController) APPAdd(c *gin.Context) {
	params := &dto.APPAddHttpInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 验证 app_id 是否被占用
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := search.Find(c, lib.GORMDefaultPool, search); err == nil {
		middleware.ResponseError(c, 2002, errors.New("租户 ID 被占用，请重新输入"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.Secret)
	}
	tx := lib.GORMDefaultPool
	info := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// APPStatistics godoc
// @Summary 租户统计
// @Description 租户统计
// @Tags 租户管理
// @ID /app/app_stat
// @Accept json
// @Product json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=StatisticsOutput} "success"
// @Router /app/app_stat [get]
func (admin *APPController) APPStatistics(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 今日流量全天小时级访问统计
	todayStat := []int64{}
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + detail.AppID)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		c.Abort()
		return
	}
	// 倒推时间，根据租户的 AppID 以及对应的时间从 redis 中获取访问数据
	currentTime := time.Now()
	for i := 0; i <= time.Now().In(lib.TimeLocation).Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayStat = append(todayStat, hourData)
	}

	// 昨日流量全天小时级访问统计
	yesterdayStat := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayStat = append(yesterdayStat, hourData)
	}
	stat := dto.StatisticsOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	}
	middleware.ResponseSuccess(c, stat)
	return
}

// APPDelete godoc
// @Summary 租户删除
// @Description 租户删除
// @Tags 租户管理
// @ID /app/app_delete
// @Accept json
// @Product json
// @Param id query string true "租户ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_delete [get]
func (admin *APPController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.IsDelete = 1
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// APPUpdate godoc
// @Summary 租户更新
// @Description 租户更新
// @Tags 租户管理
// @ID /app/app_update
// @Accept json
// @Product json
// @Param body body dto.APPAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app_update [get]
func (admin *APPController) APPUpdate(c *gin.Context) {
	params := &dto.APPUpdateHttpInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.Qps = params.Qps
	info.Qpd = params.Qpd
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}
