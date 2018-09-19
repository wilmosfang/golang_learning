package main //指明此包为 main 包

import (
	//Viper is a complete configuration solution for Go applications including 12-Factor apps. It is designed to work within an application, and can handle all types of configuration needs and formats
	//viper 是一个非常强大的用来进行配置管理的包
	"github.com/spf13/viper"
	//用来格式化输出一些内容
	"fmt"
)

//定义了一个配置函数，用来设定初始配置
//这个函数接收一个 Viper 的指针，然后对这个 Viper 结构进行配置
func confViper(v *viper.Viper) {
	//func SetConfigType(in string)
	//SetConfigType sets the type of the configuration returned by the remote source, e.g. "json".
	//这里我使用 toml 的格式来填充配置
	v.SetConfigType("toml")
	//设定配置文件的文件名，这里不要加后缀，否则找不到
	v.SetConfigName("init_conf")
	//设定找配置文件的默认默认路径
	v.AddConfigPath("./conf/") 
	// Find and read the config file
	// 读入配置文件的内容
	err := v.ReadInConfig() 
	//如果有错，就报出来
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//配置默认值，如果配置内容中没有指定，就使用以下值来作为配置值，给定默认值是一个让程序更健壮的办法
	v.SetDefault("port", "8000")
	v.SetDefault("realm", "testzone")
	v.SetDefault("key", "secret")
	v.SetDefault("tokenLookup", "header: Authorization, query: token, cookie: jwt")
	v.SetDefault("tokenHeadName", "Bearer")
	v.SetDefault("loginPath", "/login")
	v.SetDefault("authPath", "/auth")
	v.SetDefault("refreshPath", "/refresh_token")
	v.SetDefault("testPath", "/hello")
	v.SetDefault("db_host", "127.0.0.1")
	v.SetDefault("db_port", "5432")
	v.SetDefault("db_user", "postgresql")
	v.SetDefault("db_name", "testdb")
	v.SetDefault("db_password", "123456")
	v.SetDefault("admin_name", "admin")
	v.SetDefault("admin_pass", "admin")
	v.SetDefault("casbin_config", "./auth.conf")
	v.SetDefault("casbin_policy", "./auth.csv")
	v.SetDefault("logPath", "./log/")
	v.SetDefault("logFile", "gin.log")
}
