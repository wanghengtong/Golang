package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func GenerateJWT(data interface{}) string {
	// 定义 Token 的过期时间
	expirationTime := time.Now().Add(24 * time.Hour)
	// 创建 Claims
	claims := &Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名 Token
	jwtSecret := viper.GetString("auth.jwt.secret")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Println("Failed to sign token:", err)
		return ""
	}
	return signedToken
}

func ParseJWT(cookie string) (bool, error) {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// 使用相同的密钥来验证 Token
		jwtSecret := viper.GetString("auth.jwt.secret")
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 将解析后的 Claims 存储在上下文中
		logrus.Infof("Data: %v", claims["data"])
		return true, nil
	} else {
		logrus.Error("Token 解析失败")
		return false, err
	}
}

// 自定义 Claims 结构体
type Claims struct {
	Data interface{} `json:"data"`
	jwt.RegisteredClaims
}
