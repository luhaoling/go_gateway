package load_balance

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBanlanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom: // 随机负载均衡
		return &RandomBalance{}
	case LbConsistentHash: // 一致性哈希负载均衡
		return NewConsistentHashBanlance(10, nil)
	case LbRoundRobin: // 轮询负载均衡
		return &RoundRobinBalance{}
	case LbWeightRoundRobin: // 加权轮询负载均衡
		return &WeightRoundRobinBalance{}
	default:
		return &RandomBalance{}
	}
}

func LoadBanlanceFactorWithConf(lbType LbType, mConf LoadBalanceConf) LoadBalance {
	switch lbType {
	case LbRandom:
		lb := &RandomBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbConsistentHash:
		lb := NewConsistentHashBanlance(10, nil)
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	case LbWeightRoundRobin:
		lb := &WeightRoundRobinBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	default:
		lb := &RandomBalance{}
		lb.SetConf(mConf)
		mConf.Attach(lb)
		lb.Update()
		return lb
	}
}
