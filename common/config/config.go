package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hkyangyi/newe/common/file"
)

type Config struct {
	HTTP_RunMode         string //运行模式debug or release
	HTTP_Port            int    //http服务端口
	HTTP_ReadTimeout     int    //读取时间
	HTTP_WriteTimeout    int    //写入时间
	HTTP_ServeUrl        string //服务地址
	HTTP_RuntimeRootPath string //日志存储目录
	HTTP_ServeCode       string //服务器编号
	HTTP_ServeIp         string // 服务端IP

	DB_Type        string //数据库类型
	DB_User        string //用户名
	DB_Password    string //密码
	DB_Host        string //链接
	DB_Name        string //数据库名称
	DB_TablePrefix string // 前缀

	REDIS_Host        string //redis 链接
	REDIS_Password    string //redis 密码
	REDIS_MaxIdle     int    // 最大空闲
	REDIS_MaxActive   int    //最大连接数
	REDIS_IdleTimeout int    //空闲超时

	IMG_PrefixUrl string //图片url前缀
	IMG_SavePath  string //图片保存路径
	IMG_MaxSize   int    //最大图片大小
	IMG_AllowExts string //图片格式

	FILE_PrefixUrl string //文件前缀
	FILE_SavePath  string //文件保存路径
	FILE_MaxSize   int    //文件最大限制
	FILE_AllowExts string //文件格式

	SYS_SQLAUTH int //是否开启数据权限
}

var Conf = &Config{}

func ReadConfig() *Config {
	path := "assets/config/config.json"

	_, err := os.Stat(path)
	b := os.IsNotExist(err)
	if b {
		file.IsNotExistMkDir("assets/config")
		fmt.Println("初始化配置文件")
		//默认配置
		Conf.HTTP_RunMode = "debug"                  //运行模式debug or release
		Conf.HTTP_Port = 80                          //http服务端口
		Conf.HTTP_ReadTimeout = 60                   //读取时间
		Conf.HTTP_WriteTimeout = 60                  //写入时间
		Conf.HTTP_ServeUrl = "http://localhost/"     //服务地址
		Conf.HTTP_RuntimeRootPath = "assets/runtime" //日志存储目录
		Conf.HTTP_ServeCode = "A"                    //服务器编号
		Conf.DB_Type = "mysql"                       //数据链接类型
		Conf.DB_User = "root"                        //用户名
		Conf.DB_Password = "newe123"                 //数据库连接密码
		Conf.DB_Host = "127.0.0.1:3306"              //链接地址
		Conf.DB_Name = "newe"                        //数据库名
		Conf.DB_TablePrefix = ""                     //数据库数据表前缀
		Conf.REDIS_Host = "127.0.0.1:6379"           //redis连接
		Conf.REDIS_Password = ""                     //redis连接密码
		Conf.REDIS_MaxIdle = 2                       //最大空闲连接数
		Conf.REDIS_MaxActive = 10                    // #在给定时间内，允许分配的最大连接数（当为零时，没有限制）
		Conf.REDIS_IdleTimeout = 200                 // #在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制）
		Conf.IMG_PrefixUrl = ""                      //图片访问URL
		Conf.IMG_SavePath = ""                       //图片上传地址
		Conf.IMG_MaxSize = 2097152                   //#图片最大
		Conf.IMG_AllowExts = ""                      //#图片格式
		Conf.FILE_PrefixUrl = ""                     //文件访问URL前缀
		Conf.FILE_SavePath = ""                      //文件保存目录
		Conf.FILE_MaxSize = 2097152                  //#文件最大
		Conf.FILE_AllowExts = ""                     //#文件格式
		Conf.SYS_SQLAUTH = 0                         //是否开启数据权限，0不开启 1开启
		Conf.Write()

	} else {
		// 打开文件
		file, _ := os.Open("assets/config/config.json")
		// 关闭文件
		defer file.Close()

		err := json.NewDecoder(file).Decode(Conf)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(Conf)
	}

	return Conf
}

func (c *Config) Write() {

	jsonData, err := json.Marshal(c)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}
	err = os.WriteFile("assets/config/config.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Write file error:", err)
		return
	}
	fmt.Println("Struct successfully converted to JSON and written to file.")
}
