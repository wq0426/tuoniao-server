package sid

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sony/sonyflake"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func NewSid() *Sid {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
	return &Sid{sf}
}
func (s Sid) GenString() (string, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return "", err
	}
	return IntToBase62(int(id)), nil
}
func (s Sid) GenUint64() (uint64, error) {
	return s.sf.NextID()
}

// 生成6位长度的验证码
func (s Sid) GenerateCode() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	code := fmt.Sprintf("%06d", r.Intn(1000000))
	return code
}

// 使用0-9a-zA-Z生成8位长度的随机字符串
func (s Sid) IntToBase62(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// 使用0-9生成8位长度的随机字符串
func (s Sid) IntToNumberRand(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	result[0] = charset[1+seededRand.Intn(9)]
	for i := range result {
		if i == 0 {
			continue
		}
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
