package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hkyangyi/newe/common/file"
)

type Config struct {
	HTTP_RunMode         string `json:"http_run_mode"`          //运行模式debug or release
	HTTP_Port            int    `json:"http_port"`              //http服务端口
	HTTP_ReadTimeout     int    `json:"http_read_timeout"`      //读取时间
	HTTP_WriteTimeout    int    `json:"http_write_timeout"`     //写入时间
	HTTP_ServeUrl        string `json:"http_serve_url"`         //服务地址
	HTTP_RuntimeRootPath string `json:"http_runtime_root_path"` //日志存储目录
	HTTP_ServeCode       string `json:"http_serve_code"`        //服务器编号
	HTTP_ServeIp         string `json:"http_serve_ip"`          //服务端IP

	DB_Type        string `json:"db_type"`         //数据库类型
	DB_User        string `json:"db_user"`         //用户名
	DB_Password    string `json:"db_password"`     //密码
	DB_Host        string `json:"db_host"`         //链接
	DB_Name        string `json:"db_name"`         //数据库名称
	DB_TablePrefix string `json:"db_table_prefix"` //前缀

	REDIS_Host        string `json:"redis_host"`         //redis 链接
	REDIS_Password    string `json:"redis_password"`     //redis 密码
	REDIS_MaxIdle     int    `json:"redis_max_idle"`     //最大空闲
	REDIS_MaxActive   int    `json:"redis_max_active"`   //最大连接数
	REDIS_IdleTimeout int    `json:"redis_idle_timeout"` //空闲超时

	IMG_PrefixUrl string `json:"img_prefix_url"` //图片url前缀
	IMG_SavePath  string `json:"img_save_path"`  //图片保存路径
	IMG_MaxSize   int    `json:"img_max_size"`   //最大图片大小
	IMG_AllowExts string `json:"img_allow_exts"` //图片格式

	FILE_PrefixUrl string `json:"file_prefix_url"` //文件前缀
	FILE_SavePath  string `json:"file_save_path"`  //文件保存路径
	FILE_MaxSize   int    `json:"file_max_size"`   //文件最大限制
	FILE_AllowExts string `json:"file_allow_exts"` //文件格式

	SYS_SQLAUTH int `json:"sys_sql_auth"` //是否开启数据权限
}

var Conf = &Config{}

func ReadConfig() *Config {
	path := "assets/config/config.json"

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file.IsNotExistMkDir("assets/config")
		fmt.Println("初始化配置文件...")
		
		// 设置默认配置
		setDefaultConfig()
		Conf.Write()
		fmt.Println("配置文件已创建，请根据实际情况修改 assets/config/config.json")
	} else if err == nil {
		// 配置文件存在，读取配置
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("打开配置文件失败: %v", err)
			setDefaultConfig()
			return Conf
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(Conf)
		if err != nil {
			fmt.Printf("解析配置文件失败: %v", err)
			setDefaultConfig()
		}
	} else {
		fmt.Printf("检查配置文件失败: %v", err)
		setDefaultConfig()
	}

	return Conf
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	Conf.HTTP_RunMode = "debug"
	Conf.HTTP_Port = 8080
	Conf.HTTP_ReadTimeout = 60
	Conf.HTTP_WriteTimeout = 60
	Conf.HTTP_ServeUrl = "http://localhost:8080/"
	Conf.HTTP_RuntimeRootPath = "assets/runtime"
	Conf.HTTP_ServeCode = "A"
	Conf.HTTP_ServeIp = "127.0.0.1"
	
	Conf.DB_Type = "mysql"
	Conf.DB_User = "root"
	Conf.DB_Password = "newe123"
	Conf.DB_Host = "127.0.0.1:3306"
	Conf.DB_Name = "newe"
	Conf.DB_TablePrefix = ""
	
	Conf.REDIS_Host = "127.0.0.1:6379"
	Conf.REDIS_Password = ""
	Conf.REDIS_MaxIdle = 10
	Conf.REDIS_MaxActive = 100
	Conf.REDIS_IdleTimeout = 300
	
	Conf.IMG_PrefixUrl = "/images"
	Conf.IMG_SavePath = "upload/images"
	Conf.IMG_MaxSize = 2097152
	Conf.IMG_AllowExts = "jpg,jpeg,png,gif"
	
	Conf.FILE_PrefixUrl = "/files"
	Conf.FILE_SavePath = "upload/files"
	Conf.FILE_MaxSize = 5242880
	Conf.FILE_AllowExts = "pdf,doc,docx,xls,xlsx,txt"
	
	Conf.SYS_SQLAUTH = 0
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
