package auth

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tweb/global"
	"tweb/model/sys"
	"tweb/utils"
)

var (
	JwtWrapper  *jwt.GinJWTMiddleware //jwt中间件
	IdentityKey = "SC_IDENTITY"       //token存取键
	LoginKey    = "SC_LOGIN"          //login请求存取键
	CookieName  = "SC_JWT"            //存放令牌的cookie名称
)

func Init() {
	var err error

	JwtWrapper, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm: "sc auth",

		//用于在会话上下文中获取token
		//handler 可以通过 ctx.Get(IdentityKey) 获取token
		IdentityKey: IdentityKey,

		//jwt私有串
		Key: []byte("SCSoftwareStudio"),

		//会话刷新最短间隔
		Timeout: 60 * time.Minute,

		//会话最长有效期
		MaxRefresh: 12 * time.Hour,

		//校验用户名和密码并生成token
		Authenticator: func(c *gin.Context) (interface{}, error) {
			v, ok := c.Get(LoginKey)
			if !ok {
				return nil, fmt.Errorf("获取登录数据失败")
			}
			login, ok := v.(*Login)
			if !ok {
				return nil, fmt.Errorf("登录数据类型错误")
			}

			//校验账号密码
			sysUser, err := sys.SysUser{}.GetSysUserByAccount(login.Account)
			if err != nil {
				return nil, err
			}
			if !utils.VerifySaltedPassword(login.Password, sysUser.Password, sysUser.Salt) {
				return nil, fmt.Errorf("账号密码错误")
			}
			login.UserId = int64(sysUser.ID)
			login.RealName = sysUser.RealName
			token := &Token{
				UserId:   login.UserId,
				Account:  login.Account,
				RealName: login.RealName,
			}
			return token, nil
		},

		//返回登录结果
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			v, _ := c.Get(LoginKey)
			login := v.(*Login)
			login.Token = token
			login.TokenExpire = expire.Format(time.RFC3339)
			c.JSON(http.StatusOK, global.Response{Code: global.ErrCodeSuccess, Msg: "成功", Data: login})
		},

		// 返回登出结果
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, global.Response{Code: global.ErrCodeSuccess, Msg: "成功"})
		},

		// 转换Token为键值对，交由jwt一起携带
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*Token); ok {
				return jwt.MapClaims{
					"user_id": fmt.Sprintf("%d", v.UserId),
					"account": v.Account,
				}
			}
			return jwt.MapClaims{}
		},

		// 从http请求中提取Token并返回给上下文
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			token := &Token{
				Account: claims["account"].(string),
			}
			token.UserId, _ = strconv.ParseInt(claims["user_id"].(string), 0, 64)
			return token
		},

		// TODO 根据Token.Permission检查api调用是否合法
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},

		// 用户名密码错误，会话数据非法、超时，调用未授权接口等情况的处理
		Unauthorized: func(c *gin.Context, code int, message string) {
			ec := global.ErrCodePriviledge
			if strings.Contains(message, "expired") {
				ec = global.ErrCodeSessionGone
			}
			c.JSON(http.StatusOK, global.Response{
				Code: ec,
				Msg:  message,
			})
		},
		// token读取位置：cookie优先，其次为http头的Authorization
		TokenLookup: fmt.Sprintf("cookie: %s, header: Authorization", CookieName),

		// http头放置token时附带的前缀
		TokenHeadName: "SC_TOKEN",

		// 时间获取函数
		TimeFunc: time.Now,

		// 阻止客户端js获取cookie
		CookieHTTPOnly: true,

		// 存储jwt令牌的cookie名称
		CookieName: CookieName,
		SendCookie: true,
	})
	if err != nil {
		panic(fmt.Errorf("create jwt middleware failed: %s", err))
	}
}

// gin中间件，用于自动刷新token
func JwtRefresher() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		now := JwtWrapper.TimeFunc().Unix()
		expireAt := int64(claims["exp"].(float64))

		token := MustGetToken(c)
		res, _ := time.ParseDuration(fmt.Sprintf("%ds", expireAt-now))
		expireDate := time.Unix(expireAt, 0)

		if expireAt-now < 3600 {
			log.Infof("会话即将超期：用户<%s> 過期時間<%s> 剩余时间<%s>", token.Account, expireDate, res)

			// 会话即将失效，尝试刷新
			// 注意：如果会话已经持续保活超过JwtWrapper.MaxRefresh指定时间，则必须重新登录
			newToken, future, err := JwtWrapper.RefreshToken(c)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					c.JSON(http.StatusOK, global.Response{
						Code: global.ErrCodeSessionGone,
						Msg:  "登录会话已失效，请重新登录",
					})
				} else {
					c.JSON(http.StatusOK, global.Response{
						Code: global.ErrCodeInternal,
						Msg:  fmt.Sprintf("会话刷新失败：%s", err),
					})
				}
				return
			}

			log.Infof("会话刷新：用户<%s> 超期时间<%s> 剩余时间<%s>", token.Account, future, future.Sub(time.Now()))

			// 通过头部返回更新后的token
			c.Header("Authorization", JwtWrapper.TokenHeadName+" "+newToken)
			c.Header("Access-Control-Expose-Headers", newToken)
		} else {
			// 通过头部返回当前token
			c.Header("Authorization", JwtWrapper.TokenHeadName+" "+jwt.GetToken(c))
		}

		c.Next()
	}
}
