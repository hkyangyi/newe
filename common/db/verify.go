package db

// 检测唯一
func VerifyOnly(tablename string, id string, where map[string]interface{}) bool {
	var count int64

	if len(id) > 0 {
		Db.Table(tablename).Where(where).Where("id <> ?", id).Count(&count)

	} else {
		Db.Table(tablename).Where(where).Count(&count)
	}
	if count > 0 {
		return true
	}
	return false
}
