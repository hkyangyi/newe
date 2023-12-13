package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/ini.v1"
)

type Config struct {
	HTTP_RunMode         string //运行模式debug or release
	HTTP_Port            int    //http服务端口
	HTTP_ReadTimeout     int    //读取时间
	HTTP_WriteTimeout    int    //写入时间
	HTTP_ServeUrl        string //服务地址
	HTTP_RuntimeRootPath string //日志存储目录
	HTTP_ServeCode       string //服务器编号

	DB_Type        string
	DB_User        string
	DB_Password    string
	DB_Host        string
	DB_Name        string
	DB_TablePrefix string

	REDIS_Host        string
	REDIS_Password    string
	REDIS_MaxIdle     int
	REDIS_MaxActive   int
	REDIS_IdleTimeout int

	IMG_PrefixUrl string
	IMG_SavePath  string
	IMG_MaxSize   int
	IMG_AllowExts string

	FILE_PrefixUrl string
	FILE_SavePath  string
	FILE_MaxSize   int
	FILE_AllowExts string
}

func GetConfig() *Config {
	var data = Config{}
	err := data.Get()
	if err != nil {
		panic(err)
	}
	return &data
}

func (a *Config) Get() error {
	path := "assets/config/config.ini"
	_, err := os.Stat(path)
	b := os.IsNotExist(err)

	if b {
		fmt.Println("初始化配置文件")
		//默认配置
		a.HTTP_RunMode = "debug"                  //运行模式debug or release
		a.HTTP_Port = 80                          //http服务端口
		a.HTTP_ReadTimeout = 60                   //读取时间
		a.HTTP_WriteTimeout = 60                  //写入时间
		a.HTTP_ServeUrl = "http://localhost/"     //服务地址
		a.HTTP_RuntimeRootPath = "assets/runtime" //日志存储目录
		a.HTTP_ServeCode = "A"                    //服务器编号
		a.DB_Type = "mysql"                       //数据链接类型
		a.DB_User = "root"                        //用户名
		a.DB_Password = "newe123"                 //newe123
		a.DB_Host = "127.0.0.1:3306"              //链接地址
		a.DB_Name = "newe"                        //数据库名
		a.DB_TablePrefix = ""
		a.REDIS_Host = "127.0.0.1:6379"
		a.REDIS_Password = ""
		a.REDIS_MaxIdle = 2       //最大空闲连接数
		a.REDIS_MaxActive = 10    // #在给定时间内，允许分配的最大连接数（当为零时，没有限制）
		a.REDIS_IdleTimeout = 200 // #在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制）

		a.IMG_PrefixUrl = "" //图片访问URL

		a.IMG_SavePath = "" //图片上传地址

		a.IMG_MaxSize = 2097152 //#图片最大

		a.IMG_AllowExts = "" //#图片格式

		a.FILE_PrefixUrl = "" //文件访问URL前缀

		a.FILE_SavePath = "" //文件保存目录

		a.FILE_MaxSize = 2097152 //#文件最大
		a.FILE_AllowExts = ""    //#文件格式
		//datamap := Struct2Map(*a)

		err := os.MkdirAll("assets/config/", os.ModePerm)
		if err != nil && !os.IsExist(err) {
			fmt.Println(err)
			return err
		}
		err = os.WriteFile(path, []byte("#配置文件"), os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
		cfg, err := ini.Load(path)
		if err != nil {
			log.Fatalf("创建配置文件失败")
			return err
		}
		err = cfg.ReflectFrom(a)
		fmt.Println(err)
		cfg.SaveTo(path)

	} else {
		fmt.Println("读取配置文件")
		cfg, err := ini.Load(path)
		if err != nil {
			return err
		}
		err = cfg.MapTo(a)
		if err != nil {

			return err
		}
	}
	return nil
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func WriteTo(s *ini.Section, obj map[string]interface{}) {
	for k, v := range obj {
		var val string

		if reflect.TypeOf(v).Name() == "int" {
			va := v.(int)
			val = strconv.Itoa(va)
		} else {
			val = v.(string)
		}
		s.NewKey(k, val)
	}
}
