package handler

import (
	//JWT Middleware for Gin framework
	//一个 Gin 的 JWT 中间件
	"github.com/appleboy/gin-jwt"
	//Package gin implements a HTTP web framework called gin
	//这个包用来实现一个 HTTP 的 web 框架
	"github.com/gin-gonic/gin"
	
)

//定义一个 ping 的处理函数
//定义一个 ping 的处理函数
func H_ping(c *gin.Context) {
	//只作简单的反馈
	//反馈内容由此规定
	//JSON 格式反馈 pong
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//定义一函数用来处理请求
func H_hello(c *gin.Context) {
	//func ExtractClaims(c *gin.Context) jwt.MapClaims
	//ExtractClaims help to extract the JWT claims
	//用来将 Context 中的数据解析出来赋值给 claims
	//其实是解析出来了 JWT_PAYLOAD 里的内容
	/*
		func ExtractClaims(c *gin.Context) jwt.MapClaims {
		claims, exists := c.Get("JWT_PAYLOAD")
		if !exists {
			return make(jwt.MapClaims)
		}

		return claims.(jwt.MapClaims)
		}
	*/
	claims := jwt.ExtractClaims(c)
	//func (c *Context) JSON(code int, obj interface{})
	//JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
	//将内容序列化成 JSON 格式，然后放到响应 body 里，同时将 Content-Type 置为 "application/json"
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World",
	})
}

