package utils

import "time"

// 常用时间格式常量
const (
	TimeLayoutFull = "2006-01-02 15:04:05"
	TimeLayoutDate = "2006-01-02"
)

// ParseToTime 根据 layout 解析时间字符串到 time.Time（使用本地时区）。
// 当 layout 为空时使用 TimeLayoutFull 作为默认格式。
// 参数顺序为 timestr, layout（layout 可选，放在后面以符合习惯）。
func ParseToTime(timestr, layout string) (time.Time, error) {
	if layout == "" {
		layout = TimeLayoutFull
	}
	return time.ParseInLocation(layout, timestr, time.Local)
}

// TimeStringToTimestamp 将时间字符串解析为 Unix 秒时间戳。
// layout 为空时使用默认格式 "2006-01-02 15:04:05"。
// 参数顺序为 timestr, layout（layout 可选，放在后面以符合习惯）。
func TimeStringToTimestamp(timestr, layout string) (int64, error) {
	t, err := ParseToTime(timestr, layout)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// TimeStringToTimestampMs 将时间字符串解析为 Unix 毫秒时间戳。
// 参数顺序为 timestr, layout（layout 可选，放在后面以符合习惯）。
func TimeStringToTimestampMs(timestr, layout string) (int64, error) {
	t, err := ParseToTime(timestr, layout)
	if err != nil {
		return 0, err
	}
	return t.UnixNano() / 1e6, nil
}

// TimestampToTimeString 把时间戳（秒或毫秒）格式化为字符串。
// 如果 ts 看起来像毫秒（> 1e12），会按毫秒处理；否则按秒处理。
// layout 为空时使用默认格式 TimeLayoutFull。
func TimestampToTimeString(ts int64, layout string) string {
	if ts == 0 {
		return ""
	}
	if layout == "" {
		layout = TimeLayoutFull
	}
	var t time.Time
	if ts > 1e12 { // very likely milliseconds
		secs := ts / 1000
		nsec := (ts % 1000) * 1e6
		t = time.Unix(secs, nsec)
	} else {
		t = time.Unix(ts, 0)
	}
	return t.Format(layout)
}

// FormatTimestamp 是 TimestampToTimeString 的别名。
func FormatTimestamp(ts int64, layout string) string {
	return TimestampToTimeString(ts, layout)
}

// NowTimestamp 返回当前时间的 Unix 秒时间戳。
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// NowTimestampMs 返回当前时间的 Unix 毫秒时间戳。
func NowTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// NowString 返回当前时间的字符串表示，layout 为空时使用默认格式。
func NowString(layout string) string {
	if layout == "" {
		layout = TimeLayoutFull
	}
	return time.Now().Format(layout)
}

// UnixMsToTime 将毫秒时间戳转换为 time.Time
func UnixMsToTime(ms int64) time.Time {
	secs := ms / 1000
	nsec := (ms % 1000) * 1e6
	return time.Unix(secs, nsec)
}

// TimeToTimestampMs 将 time.Time 转换为毫秒时间戳
func TimeToTimestampMs(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// TimeToTimestamp 将 time.Time 转换为秒时间戳
func TimeToTimestamp(t time.Time) int64 {
	return t.Unix()
}
