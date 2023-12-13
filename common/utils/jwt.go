package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key string = "newetoken"

//生成TOken
func SetToken(uuid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["uuid"] = uuid                                                //用于在controller中确定用户
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(72)).Unix() //设置过期时间为72小时后
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//解析token
func AuthToken(tokenString string) (string, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		fmt.Println("HS256的token解析错误，err:", err, tokenString)
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("ParseHStoken:claims类型转换失败")
		return "", false
	}
	return claims["uuid"].(string), true
}
