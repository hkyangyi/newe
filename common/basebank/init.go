package base

type SysDictList struct {
	ID         string `gorm:"primary_key" json:"id"  form:"id"` //
	ParentId   string `json:"parentId" form:"parentId"`         //
	ParentName string `json:"parentName" form:"code"`           //字典名称
	Value      string `json:"value"`                            //值
	Name       string `json:"name"`                             //名称
	Status     int    `json:"status"`                           //1正常，2停用，10已删除
	Sort       int    `json:"sort"`                             //排序
}

func DictInit() {
	var items []SysDictList
	MYDB.Model(&SysDictList{}).Where("parent_id = ?", "").Find(&items)
	for _, v := range items {
		rediskey := "DICT_" + v.ParentName
		//查看子类
		var its []SysDictList
		MYDB.Model(&SysDictList{}).Where("parent_id = ?", v.ID).Find(&its)

		var dictdata = make(map[string]string, len(its))
		for i := 0; i < len(its); i++ {
			ks := its[i].Value
			dictdata[ks] = its[i].Name
		}

		SetDict(rediskey, dictdata)
	}
	// 自动创建表
	// err := MYDB.AutoMigrate(&BaseHikGateList{})
	// if err != nil {
	// 	panic("failed to migrate BaseHikGateList")
	// }

	// err = MYDB.AutoMigrate(&BaseHikGateAvList{})
	// if err != nil {
	// 	panic("failed to migrate BaseHikGateAvList")
	// }
	//手动将表放入字典
	InitWeiyiUploadConfig()

}

func SetDict(key string, data map[string]string) {
	REDIS.SetLong(key, data)
}

// 根据字典编码获取字典
func GetDict(key string) map[string]string {
	rediskey := "DICT_" + key
	var data = make(map[string]string)
	if key == "WeiyiUploadConfig" {
		return InitWeiyiUploadConfig()
	}
	//查询REDIS是否存在
	if b := REDIS.Exists(key); !b {

		var (
			fdb SysDictList
			its []SysDictList
		)
		MYDB.Model(&SysDictList{}).Where("parent_name = ?", key).Find(&fdb)
		MYDB.Model(&SysDictList{}).Where("parent_id = ?", fdb.ID).Find(&its)
		for i := 0; i < len(its); i++ {
			ks := its[i].Value
			data[ks] = its[i].Name
		}
		REDIS.SetLong(rediskey, data)
		return data
	}

	REDIS.Get(rediskey, &data)
	return data
}

// 更新字典
func UpdateDict(key string) {
	rediskey := "DICT_" + key
	if key == "WeiyiUploadConfig" {
		InitWeiyiUploadConfig()
		return
	}
	var data = make(map[string]string)
	var (
		fdb SysDictList
		its []SysDictList
	)
	MYDB.Model(&SysDictList{}).Where("parent_name = ?", key).Find(&fdb)
	MYDB.Model(&SysDictList{}).Where("parent_id = ?", fdb.ID).Find(&its)
	for i := 0; i < len(its); i++ {
		ks := its[i].Value
		data[ks] = its[i].Name
	}
	REDIS.SetLong(rediskey, data)
}

type WeiyiUploadConfig struct {
	ID         string `gorm:"primary_key" json:"id" form:"id"` //
	ApiName    string `json:"api_name"`                        //接口类型
	Code       string `json:"code"`                            //接口编码
	URL        string `json:"url"`                             //接口url
	Username   string `json:"username"`                        //接口用户名
	Password   string `json:"password"`                        //接口密码
	IsOpen     int    `json:"is_open"`                         //开启上传
	PubKeyStr  string `json:"pubKeyStr"`
	UpdateTime int64  `json:"update_time"` //更新时间
}

func InitWeiyiUploadConfig() map[string]string {
	var items []WeiyiUploadConfig
	var data = make(map[string]string)
	rediskey := "DICT_WeiyiUploadConfig"
	MYDB.Model(&WeiyiUploadConfig{}).Find(&items)
	for _, v := range items {
		data[v.Code] = v.ApiName
	}
	REDIS.SetLong(rediskey, data)
	return data
}
