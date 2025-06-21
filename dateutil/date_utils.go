package dateutil

import (
	"errors"
	"time"
)

// TimeUnit 定义日期时间差的计算单位
type TimeUnit int

// 时间单位常量定义
const (
	Millisecond TimeUnit = iota // 毫秒
	SecondUnit                  // 秒
	MinuteUnit                  // 分钟
	HourUnit                    // 小时
	DayUnit                     // 天
	WeekUnit                    // 周
	MonthUnit                   // 月
	YearUnit                    // 年
	QuarterUnit                 // 季度
)

// DateRange 日期范围生成器
// start: 起始日期时间
// end: 结束日期时间
// unit: 步进单位
type DateRange struct {
	start time.Time
	end   time.Time
	unit  TimeUnit
}

// Now 返回当前本地时间
func Now() time.Time {
	return time.Now()
}

// FormatDateTime 将时间格式化为 yyyy-MM-dd HH:mm:ss 格式
// t: 待格式化的时间
// 返回值: 格式化后的字符串
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 将时间格式化为 yyyy-MM-dd 格式
// t: 待格式化的时间
// 返回值: 格式化后的日期字符串
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTime 将时间格式化为 HH:mm:ss 格式
// t: 待格式化的时间
// 返回值: 格式化后的时间字符串
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// ParseDateTime 解析 yyyy-MM-dd HH:mm:ss 格式的字符串为时间
// s: 待解析的字符串
// 返回值: 解析后的时间和可能的错误（空输入或格式错误）
func ParseDateTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("empty input string")
	}
	return time.Parse("2006-01-02 15:04:05", s)
}

// ParseDate 解析 yyyy-MM-dd 格式的字符串为时间
// s: 待解析的字符串
// 返回值: 解析后的时间和可能的错误（空输入或格式错误）
func ParseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("empty input string")
	}
	return time.Parse("2006-01-02", s)
}

// Year 获取时间的年份
// t: 时间
// 返回值: 年份（如2023）
func Year(t time.Time) int {
	return t.Year()
}

// Month 获取时间的月份（1-12）
// t: 时间
// 返回值: 月份（1表示一月，12表示十二月）
func Month(t time.Time) int {
	return int(t.Month())
}

// Day 获取时间的日份
// t: 时间
// 返回值: 日份（1-31）
func Day(t time.Time) int {
	return t.Day()
}

// Hour 获取时间的小时（24小时制）
// t: 时间
// 返回值: 小时（0-23）
func Hour(t time.Time) int {
	return t.Hour()
}

// Minute 获取时间的分钟
// t: 时间
// 返回值: 分钟（0-59）
func Minute(t time.Time) int {
	return t.Minute()
}

// Second 获取时间的秒数
// t: 时间
// 返回值: 秒数（0-59）
func Second(t time.Time) int {
	return t.Second()
}

// IsWeekend 判断时间是否为周末（周六或周日）
// t: 时间
// 返回值: 如果是周末则为true，否则为false
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// AddYears 为时间添加指定年数
// t: 原始时间
// years: 要添加的年数（可为负数）
// 返回值: 添加后的时间
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// AddMonths 为时间添加指定月数
// t: 原始时间
// months: 要添加的月数（可为负数）
// 返回值: 添加后的时间
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddDays 为时间添加指定天数
// t: 原始时间
// days: 要添加的天数（可为负数）
// 返回值: 添加后的时间
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 为时间添加指定小时数
// t: 原始时间
// hours: 要添加的小时数（可为负数）
// 返回值: 添加后的时间
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 为时间添加指定分钟数
// t: 原始时间
// minutes: 要添加的分钟数（可为负数）
// 返回值: 添加后的时间
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 为时间添加指定秒数
// t: 原始时间
// seconds: 要添加的秒数（可为负数）
// 返回值: 添加后的时间
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// DiffYears 计算两个时间的年份差
// t1: 时间1
// t2: 时间2
// 返回值: t1年份减去t2年份的结果
func DiffYears(t1, t2 time.Time) int {
	return Year(t1) - Year(t2)
}

// DiffDays 计算两个时间的天数差
// t1: 时间1
// t2: 时间2
// 返回值: t1减去t2的天数差（可能为负数）
func DiffDays(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Hours() / 24)
}

// BeginOfSecond 获取秒级别的开始时间，即毫秒部分设置为0
// date: 日期时间
// 返回值: 秒开始时间
func BeginOfSecond(date time.Time) time.Time {
	return date.Truncate(time.Second)
}

// EndOfSecond 获取秒级别的结束时间，即毫秒设置为999
// date: 日期时间
// 返回值: 秒结束时间
func EndOfSecond(date time.Time) time.Time {
	return date.Truncate(time.Second).Add(999 * time.Millisecond)
}

// BeginOfMinute 获取某分钟的开始时间
// date: 日期时间
// 返回值: 分钟开始时间
func BeginOfMinute(date time.Time) time.Time {
	return date.Truncate(time.Minute)
}

// EndOfMinute 获取某分钟的结束时间
// date: 日期时间
// 返回值: 分钟结束时间
func EndOfMinute(date time.Time) time.Time {
	return date.Truncate(time.Minute).Add(59*time.Second + 999*time.Millisecond)
}

// BeginOfHour 获取某小时的开始时间
// date: 日期时间
// 返回值: 小时开始时间
func BeginOfHour(date time.Time) time.Time {
	return date.Truncate(time.Hour)
}

// EndOfHour 获取某小时的结束时间
// date: 日期时间
// 返回值: 小时结束时间
func EndOfHour(date time.Time) time.Time {
	return date.Truncate(time.Hour).Add(59*time.Minute + 59*time.Second + 999*time.Millisecond)
}

// BeginOfDay 获取某天的开始时间
// date: 日期时间
// 返回值: 天开始时间
func BeginOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

// EndOfDay 获取某天的结束时间
// date: 日期时间
// 返回值: 天结束时间
func EndOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 23, 59, 59, 999*1e6, date.Location())
}

// BeginOfWeek 获取某周的开始时间，默认周一为一周的第一天
// date: 日期时间
// 返回值: 周开始时间
func BeginOfWeek(date time.Time) time.Time {
	return beginOfWeek(date, true)
}

// BeginOfWeekWithMondayStart 获取某周的开始时间
// date: 日期时间
// isMondayAsFirstDay: 是否周一做为一周的第一天（false表示周日做为第一天）
// 返回值: 周开始时间
func BeginOfWeekWithMondayStart(date time.Time, isMondayAsFirstDay bool) time.Time {
	return beginOfWeek(date, isMondayAsFirstDay)
}

// EndOfWeek 获取某周的结束时间，默认周日为一周的结束
// date: 日期时间
// 返回值: 周结束时间
func EndOfWeek(date time.Time) time.Time {
	return endOfWeek(date, true)
}

// EndOfWeekWithSundayEnd 获取某周的结束时间
// date: 日期时间
// isSundayAsLastDay: 是否周日做为一周的最后一天（false表示周六做为最后一天）
// 返回值: 周结束时间
func EndOfWeekWithSundayEnd(date time.Time, isSundayAsLastDay bool) time.Time {
	return endOfWeek(date, isSundayAsLastDay)
}

// beginOfWeek 计算周开始时间的内部实现
// date: 日期时间
// isMondayAsFirstDay: 是否以周一为一周的第一天
// 返回值: 周的起始时间
func beginOfWeek(date time.Time, isMondayAsFirstDay bool) time.Time {
	currentWeekday := date.Weekday()
	var offset int

	if isMondayAsFirstDay {
		// 周一为第一天
		offset = int(currentWeekday - time.Monday)
		if offset < 0 {
			offset += 7
		}
	} else {
		// 周日为第一天
		offset = int(currentWeekday - time.Sunday)
	}

	beginDate := date.AddDate(0, 0, -offset)
	year, month, day := beginDate.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

// endOfWeek 计算周结束时间的内部实现
// date: 日期时间
// isSundayAsLastDay: 是否以周日为一周的最后一天
// 返回值: 周的结束时间
func endOfWeek(date time.Time, isSundayAsLastDay bool) time.Time {
	currentWeekday := date.Weekday()
	var offset int

	if isSundayAsLastDay {
		// 周日为最后一天
		offset = int((time.Sunday - currentWeekday + 7) % 7)
	} else {
		// 周六为最后一天
		offset = int((time.Saturday - currentWeekday + 7) % 7)
	}

	endDate := date.AddDate(0, 0, offset)
	year, month, day := endDate.Date()
	return time.Date(year, month, day, 23, 59, 59, 999*1e6, date.Location())
}

// BeginOfMonth returns the first day of the month 00:00:00
func BeginOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last day of the month 23:59:59.999999999
func EndOfMonth(t time.Time) time.Time {
	// Get first day of next month and subtract one day
	firstDayNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return EndOfDay(firstDayNextMonth.AddDate(0, 0, -1))
}

// BeginOfQuarter returns the first day of the quarter 00:00:00
func BeginOfQuarter(t time.Time) time.Time {
	month := t.Month()
	var quarterStartMonth time.Month
	switch {
	case month <= 3:
		quarterStartMonth = 1
	case month <= 6:
		quarterStartMonth = 4
	case month <= 9:
		quarterStartMonth = 7
	default:
		quarterStartMonth = 10
	}
	return time.Date(t.Year(), quarterStartMonth, 1, 0, 0, 0, 0, t.Location())
}

// EndOfQuarter returns the last day of the quarter 23:59:59.999999999
func EndOfQuarter(t time.Time) time.Time {
	month := t.Month()
	var quarterEndMonth time.Month
	switch {
	case month <= 3:
		quarterEndMonth = 3
	case month <= 6:
		quarterEndMonth = 6
	case month <= 9:
		quarterEndMonth = 9
	default:
		quarterEndMonth = 12
	}
	return EndOfMonth(time.Date(t.Year(), quarterEndMonth, 1, 0, 0, 0, 0, t.Location()))
}

// BeginOfYear 返回当年的第一天00:00:00
// t: 时间
// 返回值: 当年起始时间
func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 返回当年的最后一天23:59:59.999999999
// t: 时间
// 返回值: 当年结束时间
func EndOfYear(t time.Time) time.Time {
	return EndOfMonth(time.Date(t.Year(), 12, 1, 0, 0, 0, 0, t.Location()))
}

// Yesterday 返回昨天的开始时间
// 返回值: 昨天00:00:00
func Yesterday() time.Time {
	return BeginOfDay(Now().AddDate(0, 0, -1))
}

// Tomorrow 返回明天的开始时间
// 返回值: 明天00:00:00
func Tomorrow() time.Time {
	return BeginOfDay(Now().AddDate(0, 0, 1))
}

// LastWeek 返回一周前的当前时间
// 返回值: 一周前的时间
func LastWeek() time.Time {
	return Now().AddDate(0, 0, -7)
}

// NextWeek 返回一周后的当前时间
// 返回值: 一周后的时间
func NextWeek() time.Time {
	return Now().AddDate(0, 0, 7)
}

// LastMonth 返回一个月前的当前时间
// 返回值: 一个月前的时间
func LastMonth() time.Time {
	return Now().AddDate(0, -1, 0)
}

// NextMonth 返回一个月后的当前时间
// 返回值: 一个月后的时间
func NextMonth() time.Time {
	return Now().AddDate(0, 1, 0)
}

// OffsetMillisecond 偏移毫秒数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetMillisecond(t time.Time, offset int) time.Time {
	return t.Add(time.Duration(offset) * time.Millisecond)
}

// OffsetSecond 偏移秒数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetSecond(t time.Time, offset int) time.Time {
	return t.Add(time.Duration(offset) * time.Second)
}

// OffsetMinute 偏移分钟数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetMinute(t time.Time, offset int) time.Time {
	return t.Add(time.Duration(offset) * time.Minute)
}

// OffsetHour 偏移小时数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetHour(t time.Time, offset int) time.Time {
	return t.Add(time.Duration(offset) * time.Hour)
}

// OffsetDay 偏移天数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetDay(t time.Time, offset int) time.Time {
	return AddDays(t, offset)
}

// OffsetWeek 偏移周数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetWeek(t time.Time, offset int) time.Time {
	return AddDays(t, offset*7)
}

// OffsetMonth 偏移月数
// t: 原始时间
// offset: 偏移量(正数向未来偏移，负数向历史偏移)
// 返回值: 偏移后的时间
func OffsetMonth(t time.Time, offset int) time.Time {
	return AddMonths(t, offset)
}

// Between 计算两个日期之间的差值
// begin: 起始日期
// end: 结束日期
// unit: 差值单位
// isAbs: 是否取绝对值
func Between(begin, end time.Time, unit TimeUnit, isAbs bool) int64 {
	if begin.After(end) && !isAbs {
		return -Between(end, begin, unit, false)
	}

	var diff int64

	switch unit {
	case Millisecond:
		diff = end.Sub(begin).Milliseconds()
	case SecondUnit:
		diff = int64(end.Sub(begin).Seconds())
	case MinuteUnit:
		diff = int64(end.Sub(begin).Minutes())
	case HourUnit:
		diff = int64(end.Sub(begin).Hours())
	case DayUnit:
		diff = int64(end.Sub(begin).Hours() / 24)
	case WeekUnit:
		diff = int64(end.Sub(begin).Hours() / (24 * 7))
	case MonthUnit:
		return betweenMonth(begin, end, isAbs)
	case YearUnit:
		return betweenYear(begin, end, isAbs)
	default:
		return 0
	}

	if isAbs && diff < 0 {
		return -diff
	}
	return diff
}

// BetweenMs 计算两个日期相差的毫秒数
func BetweenMs(begin, end time.Time) int64 {
	return Between(begin, end, Millisecond, true)
}

// BetweenDay 计算两个日期相差的天数
// isReset: 是否将时间重置为当天开始
func BetweenDay(begin, end time.Time, isReset bool) int64 {
	if isReset {
		begin = BeginOfDay(begin)
		end = BeginOfDay(end)
	}
	return Between(begin, end, DayUnit, true)
}

// BetweenWeek 计算两个日期相差的周数
// isReset: 是否将时间重置为当天开始
func BetweenWeek(begin, end time.Time, isReset bool) int64 {
	if isReset {
		begin = BeginOfDay(begin)
		end = BeginOfDay(end)
	}
	return Between(begin, end, WeekUnit, true)
}

// BetweenMonth 计算两个日期相差的月数
// isReset: 是否重置时间为月初
func BetweenMonth(begin, end time.Time, isReset bool) int64 {
	return betweenMonth(begin, end, isReset)
}

// BetweenYear 计算两个日期相差的年数
// isReset: 是否重置时间为年初
func BetweenYear(begin, end time.Time, isReset bool) int64 {
	return betweenYear(begin, end, isReset)
}

// betweenMonth 计算月份差的内部实现
func betweenMonth(begin, end time.Time, isReset bool) int64 {
	if isReset {
		begin = time.Date(begin.Year(), begin.Month(), 1, 0, 0, 0, 0, begin.Location())
		end = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())
	}

	yearDiff := end.Year() - begin.Year()
	monthDiff := int(end.Month()) - int(begin.Month())
	total := int64(yearDiff*12 + monthDiff)

	// 如果没有重置且起始日期的天数大于结束日期的天数，需要调整
	if !isReset && begin.Day() > end.Day() {
		// 检查是否是同一个月
		if yearDiff == 0 && monthDiff == 0 {
			return 0
		}
		total--
	}

	if total < 0 && isReset {
		return -total
	}
	return total
}

// betweenYear 计算年份差的内部实现
func betweenYear(begin, end time.Time, isReset bool) int64 {
	if isReset {
		begin = time.Date(begin.Year(), 1, 1, 0, 0, 0, 0, begin.Location())
		end = time.Date(end.Year(), 1, 1, 0, 0, 0, 0, end.Location())
	}

	yearDiff := end.Year() - begin.Year()

	// 如果没有重置且起始日期在结束日期之后，需要调整
	if !isReset {
		if begin.Month() > end.Month() || (begin.Month() == end.Month() && begin.Day() > end.Day()) {
			yearDiff--
		}
	}

	if yearDiff < 0 && isReset {
		return int64(-yearDiff)
	}
	return int64(yearDiff)
}

// IsSameTime 判断两个时间是否完全相同（精确到纳秒）
func IsSameTime(a, b time.Time) bool {
	return a.UnixNano() == b.UnixNano()
}

// IsSameDay 判断两个日期是否为同一天
func IsSameDay(a, b time.Time) bool {
	return BeginOfDay(a).Equal(BeginOfDay(b))
}

// IsSameWeek 判断两个日期是否为同一周
// isMondayFirst 指定一周的第一天是否为周一（true: 周一, false: 周日）
func IsSameWeek(a, b time.Time) bool {
	startA := BeginOfWeek(a)
	startB := BeginOfWeek(b)
	return startA.Equal(startB)
}

// IsSameMonth 判断两个日期是否为同一月
func IsSameMonth(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month()
}

// IsLeapYear 判断是否为闰年
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// AgeOfNow 根据生日计算当前年龄
func AgeOfNow(birthDay time.Time) int {
	return age(birthDay, time.Now())
}

// AgeOfNowString 解析生日字符串并计算当前年龄
// 日期格式应为YYYY-MM-DD
func AgeOfNowString(birthDayStr string) (int, error) {
	birthDay, err := time.Parse("2006-01-02", birthDayStr)
	if err != nil {
		return 0, errors.New("解析日期失败: " + err.Error())
	}
	return AgeOfNow(birthDay), nil
}

// age 计算两个日期之间的年龄差
func age(birthDay, today time.Time) int {
	if today.Before(birthDay) {
		return 0
	}

	yearDiff := today.Year() - birthDay.Year()
	birthMonth := birthDay.Month()
	birthDayOfMonth := birthDay.Day()
	todayMonth := today.Month()
	todayDayOfMonth := today.Day()

	if todayMonth < birthMonth || (todayMonth == birthMonth && todayDayOfMonth < birthDayOfMonth) {
		yearDiff--
	}

	return yearDiff
}

// GetChineseZodiac 计算生肖（仅支持1900年及以后）
// year: 农历年份
func GetChineseZodiac(year int) string {
	if year < 1900 {
		return ""
	}

	zodiacs := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	index := (year - 1900) % 12
	if index < 0 {
		index += 12
	}
	return zodiacs[index]
}

// LengthOfYear 获取指定年份的总天数
func LengthOfYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}

// LengthOfMonth 获取指定月份的总天数
// month: 月份（1-based，1表示1月，12表示12月）
// isLeapYear: 是否为闰年（仅2月需要）
func LengthOfMonth(month int, isLeapYear bool) int {
	if month < 1 || month > 12 {
		return 0
	}

	switch month {
	case 2:
		if isLeapYear {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

// IsLastDayOfMonth 判断是否为本月最后一天
// date: 待判断的日期
// 返回值: 如果是本月最后一天则返回true，否则返回false
func IsLastDayOfMonth(date time.Time) bool {
	isLeap := IsLeapYear(date.Year())
	lastDay := LengthOfMonth(int(date.Month()), isLeap)
	return date.Day() == lastDay
}

// GetLastDayOfMonth 获取本月的最后一天
// date: 日期
// 返回值: 本月最后一天的日期数字
func GetLastDayOfMonth(date time.Time) int {
	// 原代码错误地将月份作为布尔值传入 LengthOfMonth，这里需要传入是否为闰年的判断
	isLeap := IsLeapYear(date.Year())
	return LengthOfMonth(int(date.Month()), isLeap)
}

// Range 创建日期范围生成器
// start: 起始日期时间（包括）
// end: 结束日期时间（包括）
// unit: 步进单位
// 返回值: 日期范围生成器实例
func Range(start, end time.Time, unit TimeUnit) *DateRange {
	return &DateRange{
		start: start,
		end:   end,
		unit:  unit,
	}
}

// Generate 生成日期范围内的所有日期时间点
// 返回值: 日期时间点列表
func (dr *DateRange) Generate() []time.Time {
	var result []time.Time
	current := dr.start

	for !current.After(dr.end) {
		result = append(result, current)
		current = dr.next(current)
	}

	return result
}

// next 计算下一个日期时间点(内部使用)
// current: 当前日期时间
// 返回值: 下一个日期时间点
func (dr *DateRange) next(current time.Time) time.Time {
	switch dr.unit {
	case YearUnit:
		return current.AddDate(1, 0, 0)
	case MonthUnit:
		// 处理月末日期自动调整问题
		year, month, day := current.Date()
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		// 获取下一个月的最后一天
		lastDay := time.Date(nextYear, nextMonth+1, 0, 0, 0, 0, 0, current.Location()).Day()
		// 如果当前日大于下一个月的最后一天，则使用最后一天
		if day > lastDay {
			day = lastDay
		}
		return time.Date(nextYear, nextMonth, day, current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location())
	case DayUnit:
		return current.AddDate(0, 0, 1)
	case HourUnit:
		return current.Add(time.Hour)
	case MinuteUnit:
		return current.Add(time.Minute)
	case SecondUnit:
		return current.Add(time.Second)
	default:
		return current.AddDate(0, 0, 1) // 默认按天递增
	}
}

// RangeContains 计算两个日期区间的交集
// start: 第一个日期区间
// end: 第二个日期区间
// 返回值: 两个区间共有的日期时间点
func RangeContains(start, end *DateRange) []time.Time {
	if start.unit != end.unit {
		return nil
	}
	startDates := start.Generate()
	endDates := end.Generate()

	// 使用map存储endDates以便快速查找
	endDateMap := make(map[time.Time]bool)
	for _, d := range endDates {
		endDateMap[d] = true
	}

	// 找出交集
	var intersection []time.Time
	for _, d := range startDates {
		if endDateMap[d] {
			intersection = append(intersection, d)
		}
	}

	return intersection
}

// RangeNotContains 计算两个日期区间的差集 (b - a)
// start: 第一个日期区间
// end: 第二个日期区间
// 返回值: 在b区间但不在a区间的日期时间点
func RangeNotContains(start, end *DateRange) []time.Time {
	if start.unit != end.unit {
		return nil
	}
	startDates := start.Generate()
	endDates := end.Generate()

	// 使用map存储startDates以便快速查找
	startDateMap := make(map[time.Time]bool)
	for _, d := range startDates {
		startDateMap[d] = true
	}

	// 找出差集
	var difference []time.Time
	for _, d := range endDates {
		if !startDateMap[d] {
			difference = append(difference, d)
		}
	}

	return difference
}
