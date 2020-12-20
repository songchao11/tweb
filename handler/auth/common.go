package auth

import "github.com/gin-gonic/gin"

type Token struct {
	UserId   int64
	Account  string
	RealName string
}

type Login struct {
	Token       string `json:"token"`
	TokenExpire string `json:"tokenExpire"`
	UserId      int64  `json:"userId"`
	Account     string `json:"account"`
	Password    string `json:"-"`
	RealName    string `json:"realName"`
}

func MustGetToken(c *gin.Context) *Token {
	v, ok := c.Get(IdentityKey)
	if !ok {
		panic("获取token失败")
	}
	return v.(*Token)
}
