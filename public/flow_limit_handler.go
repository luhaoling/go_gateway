package public

import (
	"sync"

	"golang.org/x/time/rate"
)

var FlowLimiterHandler *FlowLimiter

type FlowLimiter struct {
	FlowLimiterMap   map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Locker           sync.RWMutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Locker:           sync.RWMutex{},
	}
}

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

func (counter *FlowLimiter) GetLimiter(serverName string, qps float64) (*rate.Limiter, error) {
	for _, item := range counter.FlowLimiterSlice {
		if item.ServiceName == serverName {
			return item.Limiter, nil
		}
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		ServiceName: serverName,
		Limiter:     newLimiter,
	}
	counter.FlowLimiterSlice = append(counter.FlowLimiterSlice, item)
	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.FlowLimiterMap[serverName] = item
	return newLimiter, nil
}
