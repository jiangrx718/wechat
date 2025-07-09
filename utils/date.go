package utils

import (
	"fmt"
	"strconv"
	"time"
)

const (
	DEFAULT_LAYOUT_DATE_TIME    = "2006-01-02 15:04:05"
	DEFAULT_LAYOUT_DATE_TIME_1  = "2006-01-02 15:04"
	DEFAULT_LAYOUT_DATE         = "2006-01-02"
	DEFAULT_LAYOUT_DATE_YMD     = "20060102"
	DEFAULT_LAYOUT_DATE_YMD_HIS = "20060102150405"
)

func Validate(date string) bool {
	if _, err := time.ParseInLocation("2006/01/02", date, time.Local); err != nil {
		if _, err := time.ParseInLocation("2006/01", date, time.Local); err != nil {
			return false
		}
	}

	return true
}

// GetCurrentDateTime 获取当前日期时间
func GetCurrentDateTime() (dateTime string) {
	dateTime = time.Now().Format(DEFAULT_LAYOUT_DATE_TIME)
	return
}

// GetCurrentDate 获取当前日期Y-M-D
func GetCurrentDate() (dateTime string) {
	dateTime = time.Now().Format(DEFAULT_LAYOUT_DATE)
	return
}

// GetCurrentDateYMD 获取当前日期YMD
func GetCurrentDateYMD() (dateTime string) {
	dateTime = time.Now().Format(DEFAULT_LAYOUT_DATE_YMD)
	return
}

// GetCurrentDateYMDHIS 获取当前日期YMDHIS
func GetCurrentDateYMDHIS() (dateTime string) {
	dateTime = time.Now().Format(DEFAULT_LAYOUT_DATE_YMD_HIS)
	return
}

// CalculateAfterDate 几天后的日期
func CalculateAfterDate(dateInt int, days int) (result int) {
	// 待处理的日期字符串
	dateStr := strconv.Itoa(dateInt)
	// 解析日期字符串
	date, err := time.Parse(DEFAULT_LAYOUT_DATE_YMD, dateStr)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return
	}
	// 增加7天
	sevenDaysLater := date.AddDate(0, 0, days)
	// 格式化为指定格式
	result, _ = strconv.Atoi(sevenDaysLater.Format(DEFAULT_LAYOUT_DATE_YMD))
	return
}

// CalculateBeforeDate 几天前的日期
func CalculateBeforeDate(dateInt int, days int) (result string) {
	// 待处理的日期字符串
	dateStr := strconv.Itoa(dateInt)
	// 解析日期字符串
	t, err := time.Parse(DEFAULT_LAYOUT_DATE_YMD, dateStr)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return
	}
	// 计算7天前的日期
	before7Days := t.AddDate(0, 0, -days)
	// 格式化输出
	fmt.Println(before7Days.Format(DEFAULT_LAYOUT_DATE_YMD))

	result = before7Days.Format(DEFAULT_LAYOUT_DATE_YMD)
	return result
}

// GetCurrentUnixTimestamp 获取当前的时间戳
func GetCurrentUnixTimestamp() (timestamp int64) {
	timestamp = time.Now().Unix()
	return
}

// GetCurrentMilliseconds 获取当前的时间戳-毫秒
func GetCurrentMilliseconds() (timestamp int64) {
	// 获取当前时间
	now := time.Now()
	// 转换为Unix时间戳，单位为秒
	seconds := now.Unix()
	// 转换为毫秒时间戳
	milliseconds := seconds * 1000
	// 加上当前时间的纳秒部分，转换为毫秒
	milliseconds += int64(now.Nanosecond()) / 1e6
	return milliseconds
}

// GetDateToUnixTimestamp 日期时间格式转换成秒时间戳
func GetDateToUnixTimestamp(inputDateTime string) (timestamp int64) {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //指定时区
	dateTime, err := time.ParseInLocation(DEFAULT_LAYOUT_DATE_TIME, inputDateTime, TimeLocation)
	if err != nil {
		return
	}
	timestamp = dateTime.Unix()
	return
}

// GetDateToUnixNanoTimestamp 日期时间格式转换成秒时间戳
func GetDateToUnixNanoTimestamp(inputDateTime string) (timestamp int64) {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //指定时区
	dateTime, err := time.ParseInLocation(DEFAULT_LAYOUT_DATE_TIME, inputDateTime, TimeLocation)
	if err != nil {
		return
	}
	timestamp = dateTime.UnixNano()
	return
}

// GetUnixTimeToDate 时间戳转日期时间
func GetUnixTimeToDateTime(timestamp int64) (dateTime string) {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //指定时区
	dateTime = time.Unix(timestamp, 0).In(TimeLocation).Format(DEFAULT_LAYOUT_DATE_TIME)
	return
}

// GetUnixTimeToDate 时间戳转日期时间
func GetUnixTimeToDateTime1(timestamp int64) (dateTime string) {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //指定时区
	dateTime = time.Unix(timestamp, 0).In(TimeLocation).Format(DEFAULT_LAYOUT_DATE_TIME_1)
	return
}

// GetUnixTimeToDate 时间戳转日期Y-M-D
func GetUnixTimeToDate(timestamp int64) (dateTime string) {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //指定时区
	dateTime = time.Unix(timestamp, 0).In(TimeLocation).Format(DEFAULT_LAYOUT_DATE)
	return
}

// GetUnixTimeToDateYMD 时间戳转日期YMD
func GetUnixTimeToDateYMD(timestamp int64) (dateTime string) {
	TimeLocation, _ := time.LoadLocation("Asia/Shanghai") //指定时区
	dateTime = time.Unix(timestamp, 0).In(TimeLocation).Format(DEFAULT_LAYOUT_DATE_YMD)
	return
}

// 封装函数，用于 defer 调用
// defer TrackTime(time.Now(), "GetUnixTimeToDateYMD")
func TrackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s 运行耗时: %s\n", name, elapsed)
}

func GetCurrTimestamp() (timestamp int64) {
	timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	return
}

// IsValidDateTime 检查输入字符串是否符合目标格式
func IsValidDateTime(input string) bool {
	// 定义目标格式（注意：必须使用 Go 的基准时间格式）
	// 尝试解析输入字符串
	_, err := time.Parse(DEFAULT_LAYOUT_DATE_TIME, input)
	return err == nil
}

// 验证结束时间是否晚于开始时间（支持自定义格式）
func IsEndAfterStart(startStr, endStr string) (bool, error) {
	// 使用标准日期时间格式（可根据需求修改）
	const layout = "2006-01-02 15:04:05"

	// 解析开始时间
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return false, fmt.Errorf("无效的开始时间格式: %v", err)
	}

	// 解析结束时间
	end, err := time.Parse(layout, endStr)
	if err != nil {
		return false, fmt.Errorf("无效的结束时间格式: %v", err)
	}

	// 严格比较时间先后
	return end.After(start), nil
}

// ConvertDateTime 时间格式转换工具
func ConvertDateTimeToYMDHIS(inputStr string) (string, error) {
	t, err := time.Parse(DEFAULT_LAYOUT_DATE_TIME, inputStr)
	if err != nil {
		return "", fmt.Errorf("时间解析失败: %w", err)
	}
	return t.Format(DEFAULT_LAYOUT_DATE_YMD_HIS), nil
}

func ConvertDateTime(inputStr string) (string, error) {
	t, err := time.Parse(DEFAULT_LAYOUT_DATE_YMD_HIS, inputStr)
	if err != nil {
		return "", fmt.Errorf("时间解析失败: %w", err)
	}
	return t.Format(DEFAULT_LAYOUT_DATE_TIME), nil
}

// TimeRange 表示一个时间段
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// SplitByYear 将起始时间到结束时间按照“年”为单位分割
func SplitByYear(start, end time.Time) ([]TimeRange, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	var ranges []TimeRange

	current := start
	for current.Before(end) {
		var next time.Time
		year := current.Year()

		if year == end.Year() {
			// 最后一个区间直接用 end 结束
			next = end
		} else {
			// 否则取到当年最后一天的 23:00:00
			lastOfYear := time.Date(year, 12, 31, 23, 0, 0, 0, current.Location())
			next = lastOfYear
		}

		ranges = append(ranges, TimeRange{
			Start: current,
			End:   next,
		})

		// 进入下一年的 1 月 1 日 00:00:00
		current = time.Date(year+1, 1, 1, 0, 0, 0, 0, current.Location())
	}

	return ranges, nil
}
