package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	v1 "app/api/v1"
	"app/internal/common"
	"app/pkg/config"
)

type SmsService interface {
	SendSmsCodeToPhone(ctx *gin.Context, req *v1.CryptoRequest) error
	GetRequestCount(ctx *gin.Context, req *v1.CryptoRequest) (int, error)
	IncrementRequestCount(ctx *gin.Context, req *v1.CryptoRequest) error
	IsWhitePhone(ctx *gin.Context, phone string) bool
}

func NewSmsService(
	service *Service,
) SmsService {
	return &smsService{
		Service: service,
	}
}

type smsService struct {
	*Service
}

func (s *smsService) SendSmsCodeToPhone(ctx *gin.Context, req *v1.CryptoRequest) error {
	code := s.sid.GenerateCode()
	service, err := NewSMSService(
		Config{
			AccessKeyID:     *tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
			AccessKeySecret: *tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),
			SignName:        *tea.String(common.SMS_SIGN),
			TemplateCode:    *tea.String(common.SMS_TEMPLATE),
		},
	)
	if err != nil {
		log.Fatalf("Error initializing SMS service: %v", err)
	}
	err = service.SendSMS(
		req.Phone, map[string]string{
			"code": code,
		},
	)
	if err != nil {
		log.Fatalf("Error sending SMS: %v", err)
	}
	// 缓存code到redis
	err = config.Rdb.Set(ctx, common.GetCryptoKey(req.Phone, common.CODE_TYPE_LOGIN), code, time.Minute*3).Err()
	if err != nil {
		s.logger.Debug("config.Rdb.Set:", err.Error())
		return err
	}
	return nil
}

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// SMSService provides functionality to send SMS
type SMSService struct {
	client *dysmsapi.Client
	config Config
}

// NewSMSService initializes a new SMSService
func NewSMSService(cfg Config) (*SMSService, error) {
	clientConfig := &openapi.Config{
		AccessKeyId:     &cfg.AccessKeyID,
		AccessKeySecret: &cfg.AccessKeySecret,
	}
	clientConfig.Endpoint = new(string)
	*clientConfig.Endpoint = "dysmsapi.aliyuncs.com"

	client, err := dysmsapi.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMS client: %w", err)
	}

	return &SMSService{client: client, config: cfg}, nil
}

// SendSMS sends an SMS to a specific phone number
func (s *SMSService) SendSMS(phoneNumber string, params map[string]string) error {
	paramStr := "{"
	for k, v := range params {
		paramStr += fmt.Sprintf("\"%s\":\"%s\",", k, v)
	}
	paramStr = paramStr[:len(paramStr)-1] + "}"

	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  &phoneNumber,
		SignName:      &s.config.SignName,
		TemplateCode:  &s.config.TemplateCode,
		TemplateParam: &paramStr,
	}
	response, err := s.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	if response.Body.Code != nil && *response.Body.Code != "OK" {
		return fmt.Errorf("failed to send SMS: %s", *response.Body.Message)
	}
	return nil
}

func (s *smsService) GetRequestCount(ctx *gin.Context, req *v1.CryptoRequest) (int, error) {
	userInfo, err := config.Rdb.Get(ctx, common.GetTodayCodeKey(req.Phone, common.CODE_TYPE_LOGIN)).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, nil
	}
	return strconv.Atoi(userInfo)
}

func (s *smsService) IncrementRequestCount(ctx *gin.Context, req *v1.CryptoRequest) error {
	key := common.GetTodayCodeKey(req.Phone, common.CODE_TYPE_LOGIN)
	_, err := config.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		config.Rdb.Set(ctx, key, 1, 1*time.Hour)
		return nil
	} else {
		err = config.Rdb.Incr(ctx, key).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// 实现CheckSmsCode方法
func CheckSmsCode(ctx context.Context, phone, code string, bisType int) error {
	// 从redis中获取code
	val, err := config.Rdb.Get(ctx, common.GetCryptoKey(phone, bisType)).Result()
	if err != nil || len(val) == 0 {
		return errors.New("验证码已过期")
	}
	if val != code {
		return errors.New("验证码错误")
	}
	return nil
}

func (s *smsService) IsWhitePhone(ctx *gin.Context, phone string) bool {
	phones := config.ConfigInstance.GetString("access.code")
	if len(phones) == 0 {
		return false
	}
	phoneList := strings.Split(phones, ",")
	for _, p := range phoneList {
		if p == phone {
			return true
		}
	}
	return false
}
