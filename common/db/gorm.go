package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hkyangyi/newe/common/worklog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type logprint struct {
}

// 写入数据库日志
func (a logprint) Printf(s string, v ...interface{}) {
	worklog.Logio.TRACE.Printf(s, v...)
}

var (
	Db *gorm.DB
)

func GormInit(host, user, pass, name string) *gorm.DB {
	newLogger := logger.New(
		logprint{}, // io writer
		logger.Config{
			SlowThreshold:              time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // Log level
			Colorful:                  false,         // 禁用彩色打印
			IgnoreRecordNotFoundError: true,          // 忽略记录未找到错误
		},
	)
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s", 
		user, pass, host, name)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置全局禁用表自动复数
			NoLowerCase:   false,
		},
		SkipDefaultTransaction: true, // 禁用默认事务
		PrepareStmt:            true,  // 预编译SQL语句
		NowFunc: func() time.Time {
			return time.Now().Local() // 使用本地时间
		},
	})
	
	if err != nil {
		worklog.Logio.WERR(fmt.Sprintf("数据库连接失败: %v, DSN: %s", err, maskDSN(user, pass, host, name)))
		// 不直接退出，让调用方处理错误
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		worklog.Logio.WERR(fmt.Sprintf("获取数据库实例失败: %v", err))
		return nil
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(20)                     // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                    // 最大打开连接数
	sqlDB.SetConnMaxLifetime(2 * time.Hour)       // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)     // 连接最大空闲时间

	// 注册自定义回调
	db.Callback().Query().After("gorm:query").Register("dictstr", DictStr)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		worklog.Logio.WERR(fmt.Sprintf("数据库连接测试失败: %v", err))
		return nil
	}

	worklog.Logio.WTRACE("数据库连接成功")
	Db = db
	go DictInit()
	return Db
}

// maskDSN 隐藏敏感信息
func maskDSN(user, pass, host, name string) string {
	if pass != "" {
		pass = "***"
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, name)
}

func DictStr(db *gorm.DB) {
	if db.Statement.Schema != nil {
		types := db.Statement.ReflectValue.Kind()
		switch types {
		case reflect.Slice, reflect.Array:
			for _, field := range db.Statement.Schema.Fields {
				// fmt.Println(field.DataType)
				cs := field.Tag.Get("dict")
				if len(cs) > 0 {
					cso := strings.Split(cs, "_")
					if len(cso) != 2 {
						break
					}
					dict := GetDict(cso[1])

					for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
						//从字段中获取数值
						if fieldValue, isZero := field.ValueOf(context.Background(), db.Statement.ReflectValue.Index(i)); !isZero {
							//设置新的值
							if fieldo := db.Statement.Schema.LookUpField(cso[0]); fieldo != nil {
								//获取val的数据类型

								dt := field.DataType
								var key string
								switch dt {
								case "string":
									key = fieldValue.(string)
									break
								case "int":
									key = strconv.Itoa(fieldValue.(int))
									break
								default:
									key = ""
								}
								if key == "" {
									continue
								}
								var value string
								if _, ok := dict[key]; ok {
									value = dict[key]
								}
								fieldo.Set(context.Background(), db.Statement.ReflectValue.Index(i), value)
							}
						}
					}
				}
			}
		case reflect.Struct:
			for _, field := range db.Statement.Schema.Fields {
				cs := field.Tag.Get("dict")
				if len(cs) > 0 {
					cso := strings.Split(cs, "_")
					if len(cso) != 2 {
						break
					}
					dict := GetDict(cso[1])
					//从字段中获取数值
					if fieldValue, isZero := field.ValueOf(context.Background(), db.Statement.ReflectValue); !isZero {
						//设置新的值
						if fieldo := db.Statement.Schema.LookUpField(cso[0]); fieldo != nil {
							//获取val的数据类型

							dt := field.DataType
							var key string
							switch dt {
							case "string":
								key = fieldValue.(string)
								break
							case "int":
								key = strconv.Itoa(fieldValue.(int))
								break
							default:
								key = ""
							}
							if key == "" {
								continue
							}
							var value string
							if _, ok := dict[key]; ok {
								value = dict[key]
							}
							fieldo.Set(context.Background(), db.Statement.ReflectValue, value)
						}
					}

				}
			}
		}
	}
}
