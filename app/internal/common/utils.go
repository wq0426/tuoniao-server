package common

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

// 实现IsPhone(phone)bool函数，判断手机号格式是否正确
func IsPhone(phone string) bool {
	// 创建一个正则表达式对象
	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 使用正则表达式对象的 MatchString 方法来判断 phone 是否符合手机号码的格式
	return reg.MatchString(phone)
}

func IsValidCode(code string) bool {
	matched, err := regexp.MatchString(`^[0-9]{6}$`, code)
	if err != nil {
		return false
	}
	if !matched {
		return false
	}
	return true
}

func IsValidInviteCode(code string) bool {
	matched, err := regexp.MatchString(`^[0-9a-zA-Z]{8}$`, code)
	if err != nil {
		return false
	}
	if !matched {
		return false
	}
	return true
}

func GetNowDateTime() *time.Time {
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil
	}
	// 将UTC时间转换为北京时间
	currentTimeBeijing := time.Now().UTC().In(beijingLocation)
	// 格式化时间为"2006-01-02 15:04:05"格式
	return &currentTimeBeijing
}

func GetNowDateTimeStr() string {
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}
	// 将UTC时间转换为北京时间
	currentTimeBeijing := time.Now().UTC().In(beijingLocation)
	// 格式化时间为"2006-01-02 15:04:05"格式
	return currentTimeBeijing.Format(STANDARD_TIMESTR)
}

func GetNowTimestamp() int64 {
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}
	// 将UTC时间转换为北京时间
	currentTimeBeijing := time.Now().UTC().In(beijingLocation)
	return currentTimeBeijing.Unix()
}

// StringToTimestamp 将字符串转换为时间戳
func StringToTimestamp(timeStr string) (*time.Time, error) {
	// 加载北京时间位置
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	// 解析时间字符串
	t, err := time.ParseInLocation(STANDARD_TIMESTR, timeStr, beijingLocation)
	if err != nil {
		return nil, err
	}
	// 返回时间戳
	return &t, nil
}

func GetTodayCodeKey(phone string, bisType int) string {
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}
	// 将UTC时间转换为北京时间
	currentTimeBeijing := time.Now().UTC().In(beijingLocation)
	return GetCryptoKey(phone, bisType) + "." + currentTimeBeijing.Format("01-02")
}

// 根据phone和type拼接key
func GetCryptoKey(suffix string, t int) string {
	switch t {
	case CODE_TYPE_REGISTER:
		return REDIS_PREFIX_REGISTER + suffix
	case CODE_TYPE_LOGIN:
		return REDIS_PREFIX_LOGIN + suffix
	case CODE_TYPE_RESET:
		return REDIS_PREFIX_RESET + suffix
	case CODE_TYPE_CHANGE:
		return REDIS_PREFIX_CHANGE + suffix
	case CODE_TYPE_REBIND:
		return REDIS_PREFIX_REBIND + suffix
	case CODE_TYPE_FRIEND:
		return REDIS_PREFIX_FRIEND + suffix
	case CODE_TYPE_REALTIME:
		return REDIS_PREFIX_REALTIME + suffix
	case CODE_TYPE_DEVICE_UPGRADE:
		return REDIS_PREFIX_DEVICE_UPGRADE + suffix
	case CODE_TYPE_LOGIN_EXPIRE:
		return REDIS_PREFIX_LOGIN_EXPIRE + suffix
	}
	return string(t) + "." + suffix
}

// 过滤json字符串中的\n
func FilterJsonString(str string) string {
	reg := regexp.MustCompile(`\t|\n|\r`)
	return reg.ReplaceAllString(str, "")
}

func Marshal(v interface{}) string {
	str, _ := json.Marshal(v)
	return string(str)
}

func Unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// 根据ROLE_NORMAL和ROLE_VIP返回对应的角色名称
func GetRoleName(role int8) string {
	switch role {
	case ROLE_NORMAL:
		return "普通用户"
	case ROLE_VIP:
		return "vip会员"
	}
	return ""
}

// 验证交友密码格式是否对
func IsValidSocialPassword(password string) bool {
	matched, err := regexp.MatchString(`^[0-9]{6}$`, password)
	if err != nil {
		return false
	}
	if !matched {
		return false
	}
	return true
}

// IsValidURL
func IsValidURL(url string) bool {
	matched, err := regexp.MatchString(`^http[s]?://.*`, url)
	if err != nil {
		return false
	}
	if !matched {
		return false
	}
	return true
}

// CalculateProgressPercentage 计算当前时间距离开始时间占总时长的比例
func CalculateProgressPercentage(startTime, endTime, curTime time.Time) int {
	totalDuration := endTime.Sub(startTime).Seconds()
	elapsedDuration := curTime.Sub(startTime).Seconds()
	progressPercentage := int((elapsedDuration / totalDuration) * 100)
	return progressPercentage
}

// 保留6位小数
func FormatFloat(f float64) float64 {
	formattedNumStr := fmt.Sprintf("%.3f", f)
	formattedNum, err := strconv.ParseFloat(formattedNumStr, 64)
	if err != nil {
		return 0
	}
	return formattedNum
}

// 将浮点数向下取整
func FloorToInt(f float64) int {
	return int(math.Floor(f))
}

// 获取escort_record表对应的字段名
func GetEscortFieldName(rewardType int8) (string, string) {
	switch rewardType {
	case REWARD_TYPE_HORSE:
		return "today_horse", "total_horse"
	case REWARD_TYPE_GRAIN:
		return "today_grain", "total_grain"
	case REWARD_TYPE_IRON:
		return "today_iron", "total_iron"
	case REWARD_TYPE_PLUNDER:
		return "today_plunder", "total_plunder"
	case REWARD_TYPE_PUTDOWN:
		return "today_putdown", "total_putdown"
	case REWARD_TYPE_BATTLE:
		return "today_battle", "total_battle"
	case REWARD_TYPE_ESCORT:
		return "today_escort", "total_escort"
	case REWARD_TYPE_ADVERTISING:
		return "today_advertising", "total_advertising"
	case REWARD_TYPE_DONATE:
		return "today_donate", "total_donate"
	}
	return "", ""
}

// 获取account表对应的字段名
func GetAccountFieldName(rewardType int8) []string {
	switch rewardType {
	case REWARD_TYPE_HORSE:
		return []string{"horse"}
	case REWARD_TYPE_GRAIN:
		return []string{"grain"}
	case REWARD_TYPE_IRON:
		return []string{"iron"}
	case REWARD_TYPE_PLUNDER:
		return []string{"horse", "grain", "iron"}
	case REWARD_TYPE_PUTDOWN:
		return []string{"prestige"}
	case REWARD_TYPE_BATTLE:
		return []string{"silver"}
	case REWARD_TYPE_ESCORT:
		return []string{"horse"}
	case REWARD_TYPE_ADVERTISING:
		return []string{}
	case REWARD_TYPE_DONATE:
		return []string{"prestige"}
	}
	return []string{}
}

// 将字符串转成整形
func StrToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

// 将字符串转成浮点数
func StrToFloat(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return num
}

func RandomKeyByWeight(weights map[int]int) int {
	// Calculate the total weight
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	// Generate a random number between 0 and totalWeight
	src := rand.NewSource(time.Now().UnixNano())
	randomWeight := rand.New(src).Intn(totalWeight)

	// Find the key corresponding to the random weight
	tmpRand := weights[0]
	for key, weight := range weights {
		if randomWeight < tmpRand {
			return key
		}
		tmpRand += weight
	}

	return -1 // This should never happen if the input map is not empty
}

func RandomInitTwoNum(max int) []int {
	var num1, num2 int
	for {
		src := rand.NewSource(time.Now().UnixNano())
		num1 = rand.New(src).Intn(max) + 1
		src2 := rand.NewSource(time.Now().UnixNano())
		num2 = rand.New(src2).Intn(max) + 1
		if num1 != num2 {
			break
		}
	}
	return []int{num1, num2}
}

func GetNowDateStr() string {
	nowTime := GetNowDateTime()
	if nowTime == nil {
		return ""
	}
	// 格式化时间为"2006-01-02 15:04:05"格式
	return nowTime.Format(DATE_FORMAT)
}

// 生成8位有a-zA-Z组成的随机字符串
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 生成订单号
func GenerateOrderNo() string {
	return "T" + time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000000))
}

// 验证手机号
func IsMobile(phone string) bool {
	matched, err := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	if err != nil {
		return false
	}
	return matched
}

// 验证身份证号
func IsIDCard(idCard string) bool {
	matched, err := regexp.MatchString(`^[0-9]{17}[0-9xX]$`, idCard)
	if err != nil {
		return false
	}
	return matched
}
