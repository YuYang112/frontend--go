package authMiddleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 签名需要传递的参数
type HmacUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MyClaims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 登录的参数
type LoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 证书签名密钥
var jwtKey = []byte("abc")

// 定义解析token的方法
func parseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// 定义中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authorization := c.GetHeader("authorization")
		if authorization == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "必须传递token",
			})
			c.Abort()
			return
		}
		tokeStrings := strings.Split(c.GetHeader("authorization")," ")
		if len(tokeStrings)!=2{
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "token格式不正确",
			})
			c.Abort()
			return
		}
		tokeString := tokeStrings[1]
		fmt.Println(tokeString, "当前token")
		token, claims, err := parseToken(tokeString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "token解析错误",
			})
			c.Abort()
			return
		}
		// 从token中解析出来的数据挂载到上下文上,方便后面的控制器使用
		c.Set("email", claims.Email)
		c.Set("password", claims.Password)
		c.Next()
	}
}

// 定义生成token的方法
func GenerateToken(u HmacUser) (string, error) {
	// 定义过期时间,7天后过期
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &MyClaims{
		Email:   u.Email,
		Password: u.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // 发布时间
			Subject:   "token",               // 主题
			Issuer:    "水痕",                  // 发布者
		},
	}
	// 注意单词别写错了
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
