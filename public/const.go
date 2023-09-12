package public

const (
	ValidatorKey        = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	AdminSessionInfoKey = "AdminSessionInfoKey"

	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0 // url负载均衡(使用请求的前缀 url 进行流量路由)
	HTTPRuleTypeDomain    = 1 // 域名负载均衡(使用域名来进行流量路由)

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"

	JwtSignKey = "my_sign_key"
	JwtExpired = 60 * 60 // jwt 过期时间 1h
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)
