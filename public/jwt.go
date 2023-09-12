package public

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	//"github.com/dgrijalva/jwt-go"
)

func JwtDecode(tokenString string) (*jwt.StandardClaims, error) {
	// 将令牌字符串解析为一个 jwt.Token 对象，并且使用 JwtSignKey 作为密钥进行验证
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("token is not jwt.StandardClaims")
	}
}

func JwtEncode(claims jwt.StandardClaims) (string, error) {
	mySigningKey := []byte(JwtSignKey)
	// 使用 jwt.SigningMethodHS256 算法进行签名，使用传入的 claims 最为有效载荷
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用 mySigningKey 对令牌进行签名
	return token.SignedString(mySigningKey)
}
