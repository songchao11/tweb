package auth

import "github.com/gin-gonic/gin"

type Token struct {
	UserId   int64
	Username string
}

type Login struct {
	Token       string `json:"token"`
	TokenExpire string `json:"tokenExpire"`
	UserId      int64  `json:"userId"`
	Username    string `json:"username"`
	Password    string `json:"-"`
}

func MustGetToken(c *gin.Context) *Token {
	v, ok := c.Get(IdentityKey)
	if !ok {
		panic("获取token失败")
	}
	return v.(*Token)
}
