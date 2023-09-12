package public

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
)

// 生成一个含有盐值的密码字符串
func GenSaltPassword(salt, password string) string {
	// 创建一个 sha256 的哈希算法实例
	s1 := sha256.New()
	// 密码转化为字符数组
	s1.Write([]byte(password))
	// 使用 s1 进行哈希运算，并转化为字符串
	str1 := fmt.Sprintf("%x", s1.Sum(nil))

	// 创建另外一个 sha256 哈希算法，并且将 str1 和 salt 连接起来，转换为字符串，并且使用 s2 进行哈希运算
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}

// MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Obj2Json(s interface{}) string {
	bts, _ := json.Marshal(s)
	return string(bts)
}

func InStringSlice(slice []string, str string) bool {
	for _, item := range slice {
		if str == item {
			return true
		}
	}
	return false
}
