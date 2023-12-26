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
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // 禁用彩色打印
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			//设置全局禁用表自动复数
			SingularTable: true,
		},
		SkipDefaultTransaction: true, //禁用默认事务
	})
	if err != nil {
		log.Fatalf("model.Setup err: %v", err)
	}
	db.Callback().Query().After("gorm:query").Register("dictstr", DictStr)

	sqlDB, err := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	Db = db.Debug()
	return Db
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
