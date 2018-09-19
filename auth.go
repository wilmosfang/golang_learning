package main //指明此包为 main 包


import (
	//Package time provides functionality for measuring and displaying time
	//The calendrical calculations always assume a Gregorian calendar, with no leap seconds
	//用来提供一些计算和显示时间的函数
	"time"
	//JWT Middleware for Gin framework
	//一个 Gin 的 JWT 中间件
	"github.com/appleboy/gin-jwt"
	//Package gin implements a HTTP web framework called gin
	//这个包用来实现一个 HTTP 的 web 框架
	"github.com/gin-gonic/gin"
	//加入postgres的库
	_ "github.com/jinzhu/gorm/dialects/postgres"
)



//定义一个回调函数，用来决断用户id和密码是否有效
func authCallback( c *gin.Context) (interface{}, error) {
	//userId string, password string,
	var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
	uid := loginVals.Username
	pw := loginVals.Password
	//这里的通过从数据库中查询来判断是否为现存用户，生产环境下一般都会使用数据库来存储账号信息，进行检验和判断
	user := User{} //创建一个临时的存放空间
	//如果这条记录存在的的情况下
	if !authDB.Where("user_id = ?", uid).Find(&user).RecordNotFound() {
		//定义一个临时的结构对象
		queryRes := User{} //创建一个临时的存放空间
		//将 user_id 为认证信息中的 密码找出来(目前密码是明文的，这个其实不安全，可以通过加盐哈希将结果进行对比的方式以提高安全等级，这里只作原理演示，就不搞那么复杂了)
		//找到后放到前面定义的临时结构变量里
		authDB.Where("user_id = ?", uid).Find(&queryRes)
		//对比，如果密码也相同，就代表认证成功了
		if queryRes.Password == pw {
			//反馈相关信息和 true 的值，代表成功
			return &login{
				Username: uid,
				Password: pw,
			}, nil
		}
	}
	//否则返回失败
	return nil, jwt.ErrFailedAuthentication
}

//定义一个回调函数，用来决断用户在认证成功的前提下，是否有权限对资源进行访问
func authPrivCallback(data interface{}, c *gin.Context) bool {
	claims := jwt.ExtractClaims(c)
	return casbinEnforcer.Enforce(claims["id"], c.Request.URL.String(), c.Request.Method)
}

//定义一个函数用来处理，认证不成功的情况
func unAuthFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func payLoad(data interface{})  jwt.MapClaims {
	if v, ok := data.(*login); ok {
		return jwt.MapClaims{
			"id": v.Username,
		}
	}
	return jwt.MapClaims{}
}



//定义一个中间件，用来反馈 jwt 的认证逻辑
//这里将相应的配置直接以变量的方式传递进来了
func addMidd(v_realm, v_key, v_tokenLookup, v_tokenHeadName string) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		//Realm name to display to the user. Required.
		//必要项，显示给用户看的域
		Realm: v_realm,
		//Secret key used for signing. Required.
		//用来进行签名的密钥，就是加盐用的
		Key: []byte(v_key),
		//Duration that a jwt token is valid. Optional, defaults to one hour
		//JWT 的有效时间，默认为一小时
		Timeout: time.Hour,
		// This field allows clients to refresh their token until MaxRefresh has passed.
		// Note that clients can refresh their token in the last moment of MaxRefresh.
		// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
		// Optional, defaults to 0 meaning not refreshable.
		//最长的刷新时间，用来给客户端自己刷新 token 用的
		//IdentityKey: identityKey,
	
		
		PayloadFunc: payLoad,
		
		MaxRefresh: time.Hour,
		// Callback function that should perform the authentication of the user based on userID and
		// password. Must return true on success, false on failure. Required.
		// Option return user data, if so, user data will be stored in Claim Array.
		//必要项, 这个函数用来判断 User 信息是否合法，如果合法就反馈 true，否则就是 false, 认证的逻辑就在这里
		Authenticator: authCallback,
		// Callback function that should perform the authorization of the authenticated user. Called
		// only after an authentication success. Must return true on success, false on failure.
		// Optional, default to success
		//可选项，用来在 Authenticator 认证成功的基础上进一步的检验用户是否有权限，默认为 success
		Authorizator: authPrivCallback,
		// User can define own Unauthorized func.
		//可以用来息定义如果认证不成功的的处理函数
		Unauthorized: unAuthFunc,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		//这个变量定义了从请求中解析 token 的位置和格式
		TokenLookup: v_tokenLookup,
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		//TokenHeadName 是一个头部信息中的字符串
		TokenHeadName: v_tokenHeadName,
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		//这个指定了提供当前时间的函数，也可以自定义
		TimeFunc: time.Now,
	}
}

