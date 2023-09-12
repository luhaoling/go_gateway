package load_balance

type LoadBalance interface {
	Add(...string) error        // 添加
	Get(string) (string, error) // 获取
	//后期服务发现补充
	Update() // 更新
}
