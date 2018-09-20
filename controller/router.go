package controller //指明此包为 main 包

import (
	//os 包可以提供抽象的操作系统函数
	"os"
	//io 包可以提供基础的接口以进行I/O类操作
	"io"
	//Package bytes implements functions for the manipulation of byte slices. It is analogous to the facilities of the strings package
	//这个包用来提供对 byte slices 的操作，有点类似于 string 的作用
	//"bytes"
	//用来格式化输出一些内容
	"fmt"
	//Package time provides functionality for measuring and displaying time
	//The calendrical calculations always assume a Gregorian calendar, with no leap seconds
	//用来提供一些计算和显示时间的函数
	//"time"
	//JWT Middleware for Gin framework
	//一个 Gin 的 JWT 中间件
	//"github.com/appleboy/gin-jwt"
	//Package gin implements a HTTP web framework called gin
	//这个包用来实现一个 HTTP 的 web 框架
	"github.com/gin-gonic/gin"
	//Viper is a complete configuration solution for Go applications including 12-Factor apps. It is designed to work within an application, and can handle all types of configuration needs and formats
	//viper 是一个非常强大的用来进行配置管理的包
	"github.com/spf13/viper"
	//The fantastic ORM library for Golang
	//一个非常好用的 go 语言 ORM 库
	"github.com/jinzhu/gorm"
	//加入postgres的库
	_ "github.com/jinzhu/gorm/dialects/postgres"
	//An authorization library that supports access control models like ACL, RBAC, ABAC in Golang
	//一个支持 ACL, RBAC, ABAC 访问控制模型的 go 语言库
	"github.com/casbin/casbin"
	//Gorm Adapter is the Gorm adapter for Casbin. With this library, Casbin can load policy from Gorm supported database or save policy to it.
	//可以使用这个包来完成策略在数据库中的存取
	"github.com/casbin/gorm-adapter"
	//Authz is an authorization middleware for Gin, it's based on https://github.com/casbin/casbin
	//Authz 是 Gin 的认证中间件，它是基于 casbin 实现的
	//"github.com/gin-contrib/authz"
	//"reflect"
	"gw/handler"
	"gw/model"
)

//定义一个内部全局的 db 指针用来进行认证，数据校验
var authDB *gorm.DB

//定义一个内部全局的 viper 指针用来进行配置读取
var config *viper.Viper

//定义一个内部全局的 casbin.Enforcer 指针用来进行权限校验
var casbinEnforcer *casbin.Enforcer

//main 包必要一个 main 函数，作为起点
func Init() {
	//func New() *Viper
	//New returns an initialized Viper instance.
	//用来生成一个新的 viper
	v := viper.New()
	//将全局的 config 进行赋值
	config = v
	//对 viper 进行配置
	confViper(v)
	//构建一个 pg 的连接串
	pg_conn_info := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", v.GetString("db_host"), v.GetString("db_port"), v.GetString("db_user"), v.GetString("db_name"), v.GetString("db_password"))
	db, err := gorm.Open("postgres", pg_conn_info)
	if err != nil {
		panic(err) //如果出错，就直接打印出错信息，并且退出
	} else {
		fmt.Println("Successfully connected!") //如果没有出错，就打印成功连接的信息
		authDB = db                            //连接成功的情况下将认证的数据库进行赋值
	}
	//func (s *DB) Close() error
	//Close close current db connection. If database connection is not an io.Closer, returns an error.
	//panic 之后并不是直接退出，而是先去执行 defer 的内容
	//关闭当前的 db 连接
	defer db.Close() //如果出错，先将 db 关掉
	//func (s *DB) AutoMigrate(values ...interface{}) *DB
	//AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
	//AutoMigrate 会自动将给定的模型进行迁移，只会添加缺失的字段，并不会删除或者修改当前的字段
	db.AutoMigrate(&model.User{})
	//创建一个结构变量
	user := model.User{}
	//如果 db 中没有这条记录，就创建，如果有就忽略掉
	if db.Where("user_id = ?", v.GetString("admin_name")).Find(&user).RecordNotFound() {
		user := model.User{UserID: v.GetString("admin_name"), Password: v.GetString("admin_pass")}
		db.Create(&user)
	}
	//这里有一个坑，是通过看源码解决的
	//如果不手动指定，而是自动创建时，它会默认首先要求一个很大的 postgresql 权限，尝试在 postgres 下面创建一个库，如果过程报错，就直接 panic 出来
	//这里我就通过手动创建的方式，直接指到 testdb 里
	//它会自动创建一个叫 casbin_rule 的表来进行规则存放
	casbin_adapter := gormadapter.NewAdapter("postgres", pg_conn_info, true)
	//使用前面定的 casbin_adapter 来构建 enforcer
	e := casbin.NewEnforcer(v.GetString("casbin_config"), casbin_adapter)
	//将全局的 enforcer 进行赋值，以方便在其它地方进行调用
	casbinEnforcer = e
	//加载规则
	e.LoadPolicy()
	//下面的这些命令可以用来添加规则
	e.AddPolicy(v.GetString("admin_name"), v.GetString("authPath")+v.GetString("testPath"), "GET")
	//e.AddPolicy("dex", "/auth/hello", "GET")
	//e.AddRoleForUser("user_a", "user")
	//e.AddRoleForUser("user_b", "user")
	//e.AddRoleForUser("user_c", "user")
	e.AddPolicy("user", v.GetString("authPath")+"/ping", "GET")
	e.SavePolicy()

	f, _ := os.Create(v.GetString("logPath") + v.GetString("logFile"))                     //func Create(name string) (*File, error) 接受一个文件名字符串，反馈一个文件指针，和一个错误输出
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //将输出写出到文件与终端各一份
	//func Default() *Engine
	//Default returns an Engine instance with the Logger and Recovery middleware already attached.
	//用来返回一个已经加载了Logger and Recovery中间件的引擎
	r := gin.Default()
	/*gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}*/
	//the jwt middleware
	authMiddleware := addMidd(v.GetString("realm"), v.GetString("key"), v.GetString("tokenLookup"), v.GetString("tokenHeadName"))

	//func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context)
	//LoginHandler can be used by clients to get a jwt token. Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}. Reply will be of the form {"token": "TOKEN"}.
	//将 /login 交给 authMiddleware.LoginHandler 函数来处理
	r.POST(v.GetString("loginPath"), authMiddleware.LoginHandler)
	//func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup
	//Group creates a new router group. You should add all the routes that have common middlwares or the same path prefix. For example, all the routes that use a common middlware for authorization could be grouped
	//创建一个组 auth
	auth := r.Group(v.GetString("authPath"))
	//func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc
	//MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
	//auth 组中使用 MiddlewareFunc 中间件
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		//如果是 /auth 组下的 /hello 就交给 helloHandler 来处理
		auth.GET(v.GetString("testPath"), handler.H_hello)
		//func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context)
		//RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh. Shall be put under an endpoint that is using the GinJWTMiddleware. Reply will be of the form {"token": "TOKEN"}.
		//如果是 /auth 组下的 /refresh_token 就交给 RefreshHandler 来处理
		auth.GET(v.GetString("refreshPath"), authMiddleware.RefreshHandler)
		auth.GET("/ping", handler.H_ping)
	}

	r.Run(":" + v.GetString("port")) //在 0.0.0.0:配置端口 上启监听
}
