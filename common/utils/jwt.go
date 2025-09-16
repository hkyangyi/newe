package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const key string = "github.com/hkyangyi/newetoken"

// 生成Token
func SetToken(uuid string) (string, error) {
	claims := jwt.MapClaims{
		"uuid": uuid, //用于在controller中确定用户
		"exp":  time.Now().Add(time.Hour * 72).Unix(), //设置过期时间为72小时后
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析token
func AuthToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", errors.New("HS256的token解析错误")
	}

	if !token.Valid {
		return "", errors.New("token无效")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("ParseHStoken:claims类型转换失败")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", errors.New("无效的过期时间")
	}

	if int64(exp) < time.Now().Unix() {
		return "", errors.New("token已失效")
	}

	uuid, ok := claims["uuid"].(string)
	if !ok {
		return "", errors.New("无效的用户标识")
	}

	return uuid, nil
}
