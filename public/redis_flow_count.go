package public

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/e421083458/go_gateway/golang_common/lib"
	"github.com/garyburd/redigo/redis"
)

type RedisFlowCountService struct {
	AppID       string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appID,
		Interval: interval,
		QPS:      0,
		Unix:     0, // 上次更新统计数据的时间戳
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		// 定义一个周期性触发定时的计时器，按照 1s 的间隔往 channel 发送系统当前时间
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C // 先阻塞 1s , 1s 后从计数器通道上接收触发事件的时间值
			// 获取当前触发周期的请求数量（从计数器通道上接收到触发事件的事件时，代表这一周期已经结束）
			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount) // 原子性获取数据（相当于 1s 内的请求数据）
			// 准备下一个触发周
			atomic.StoreInt64(&reqCounter.TickerCount, 0) // 原子性重置数据

			// 获取当前时间
			currentTime := time.Now()
			// 建立以天为单位，以小时为单位的统计指标
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)
			// 封装事务
			if err := RedisConfPipline(func(c redis.Conn) {
				c.Send("INCRBY", dayKey, tickerCount) // 将 dayKey 对应的数值加上 tickerCount
				c.Send("EXPIRE", dayKey, 86400*2)     // 86400*2 两天的时间 将 daykey 的生存时间设置为 2 天
				c.Send("INCRBY", hourKey, tickerCount)
				c.Send("EXPIRE", hourKey, 86400*2)
			}); err != nil {
				fmt.Println("RedisConfPipline err:", err)
				continue
			}
			// 获取当前时间的日统计数据
			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				fmt.Println("reqCounter.GetDayData err:", err)
				continue
			}
			// 确定时间
			nowUnix := time.Now().Unix() // 当前本地事件的 UNIX 时间戳
			// 如果此前没有更新过统计事件的事件戳，则设置 reqCounter.Unix 的时间戳
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue //
			}
			// 计算本周期内的请求数量
			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()
	return reqCounter
}

func (o *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(lib.TimeLocation).Format("20060102")              // 精确到日（当前时间）
	return fmt.Sprintf("%s_%s_%s", RedisFlowDayKey, dayStr, o.AppID) // 日前缀+当前日+服务名称（全站 服务 租户）
}

func (o *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocation).Format("2006010215")             // 精确到时（当前时间）
	return fmt.Sprintf("%s_%s_%s", RedisFlowHourKey, hourStr, o.AppID) // 时前缀+当前时+服务名称（全站 服务 租户）
}

func (o *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", o.GetHourKey(t)))
}

func (o *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", o.GetDayKey(t)))
}

// 原子增加
func (o *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("err")
			}
		}()
		atomic.AddInt64(&o.TickerCount, 1)
	}()
}
