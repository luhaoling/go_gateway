package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	//default check setting
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

type LoadBalanceCheckConf struct {
	observers    []Observer
	confIpWeight map[string]string // ip-权重 map
	activeList   []string          // 活跃的 ip 列表
	format       string            // ip 和 权重的输出形式
}

func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// 获取 conf
func (s *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50" //默认weight
		}
		// conf 的输出形式
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

//更新配置时，通知监听者也更新
func (s *LoadBalanceCheckConf) WatchConf() {
	//fmt.Println("watchConf")
	go func() {
		// 记录每个 ip 地址对应的错误次数
		confIpErrNum := map[string]int{}
		for {
			// 存储发生变化的 ip 地址列表
			changedList := []string{}
			// 主动探测下游节点，判断下游节点是否在线
			for item, _ := range s.confIpWeight {
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				//todo http statuscode
				// 在线
				if err == nil {
					//关闭连接
					conn.Close()
					// 将 ip 列表的错误计数置为0
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}
				// 不在线
				if err != nil {
					// ip 列表的错误计数加1
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] += 1
					} else {
						confIpErrNum[item] = 1
					}
				}
				// 如果错误计数小于 2，才能将 ip 地址添加到 changedList 列表中
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			// 对两个列表进行排序
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			// 如果两个列表不相等，则说明配置发生了变化（下游服务器发生了变更），更新配置
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}
			// 休眠 5s 后再继续进行探测
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

//更新配置时，通知监听者也更新
func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	//fmt.Println("UpdateConf", conf)
	// 更新下游服务 ip 列表
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceCheckConf(format string, conf map[string]string) (*LoadBalanceCheckConf, error) {
	aList := []string{} // IP 列表
	//默认初始化
	for item, _ := range conf {
		aList = append(aList, item)
	}
	mConf := &LoadBalanceCheckConf{format: format, activeList: aList, confIpWeight: conf}
	mConf.WatchConf()
	return mConf, nil
}
