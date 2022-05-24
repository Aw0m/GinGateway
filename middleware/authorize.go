package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

var TokenSalt = []byte("dev123")

type MyCustomClaims struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		//log.Println("Authorization Token: ", tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return TokenSalt, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			log.Println("没有访问token，拦截", err)
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			fmt.Printf("Username: %v, ExpiresAt %v   ", claims.UserName, claims.StandardClaims.ExpiresAt)
			// 验证通过，会继续访问下一个中间件
			c.Next()
		} else {
			fmt.Println("出现错误：", err)
			// 验证不通过，不再调用后续的函数处理
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			log.Println("token错误，拦截")
			c.Abort()
		}
	}
}

func GetOpenid(tokenString string) string {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TokenSalt, nil
	})
	if err != nil {
		return ""
	}
	claims := token.Claims.(*MyCustomClaims)
	openid := claims.UserId
	return openid
}
