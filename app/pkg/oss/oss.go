// @program:     flashbear
// @file:        oss.go.go
// @author:      ac
// @create:      2024-10-28 12:15
// @description:
package oss

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"

	"app/internal/common"
)

type Oss struct {
	Base64   string `json:"base64"`
	FilePath string `json:"file_path"`
	Ext      string `json:"ext"`
}

func NewOss(filePath, ext string) *Oss {
	return &Oss{
		FilePath: filePath,
		Ext:      ext,
	}
}

// 上传文件
func (o *Oss) UploadFile(ctx *gin.Context) (string, error) {
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Fatalf("Failed to create credentials provider: %v", err)
		return "", err
	}
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	clientOptions = append(clientOptions, oss.Region(common.REGION_CN_BEIJING))
	// 设置签名版本
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	client, err := oss.New(common.HTTPS_PREFIX+common.OSS_ENDPOINT, os.Getenv("OSS_ACCESS_KEY_ID"),
		os.Getenv("OSS_ACCESS_KEY_SECRET"))
	if err != nil {
		log.Fatalf("Failed to create OSS client: %v", err)
		return "", err
	}
	bucket, err := client.Bucket(common.BUCKET_NAME)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
		return "", err
	}
	// 获取filePath中的文件名以及后缀
	fielname := filepath.Base(o.FilePath) + o.Ext
	src := rand.NewSource(time.Now().UnixNano())
	randInt := rand.New(src).Intn(10)
	objectKey := common.BUCKET_NAME + "/" + strconv.Itoa(randInt) + "/" + fielname // 请替换为实际的对象Key
	err = bucket.PutObjectFromFile(objectKey, o.FilePath)
	if err != nil {
		log.Fatalf("Failed to put object from file: %v, filePath: %s", err, o.FilePath)
		return "", err
	}
	signedURL := fmt.Sprintf("%s/%s", common.RESOURCE_AVATAR, objectKey)
	return signedURL, nil
}

func FullUrlWithDomain(ctx *gin.Context, url string) string {
	if url == "" {
		return ""
	}
	return common.HTTP_PREFIX + ctx.Request.Host + "/" + url
}

func NewOssBase64(base64 string) *Oss {
	return &Oss{
		Base64: base64,
	}
}

// base64上传文件
func (o *Oss) UploadBase64(ctx *gin.Context) (string, error) {
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Fatalf("Failed to create credentials provider: %v", err)
		return "", err
	}
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	clientOptions = append(clientOptions, oss.Region(common.REGION_CN_BEIJING))
	// 设置签名版本
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	client, err := oss.New(common.HTTPS_PREFIX+common.OSS_ENDPOINT, os.Getenv("OSS_ACCESS_KEY_ID"),
		os.Getenv("OSS_ACCESS_KEY_SECRET"))
	if err != nil {
		log.Fatalf("Failed to create OSS client: %v", err)
		return "", err
	}
	bucket, err := client.Bucket(common.BUCKET_NAME)
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
		return "", err
	}
	src := rand.NewSource(time.Now().UnixNano())
	randInt := rand.New(src).Intn(10)
	objectKey := common.BUCKET_NAME + "/" + strconv.Itoa(randInt) + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".png" // 请替换为实际的对象Key

	// 解码Base64
	decodedData, err := base64.StdEncoding.DecodeString(o.Base64)
	if err != nil {
		log.Fatalf("Failed to decode base64: %v", err)
		return "", err
	}

	err = bucket.PutObject(objectKey, bytes.NewReader(decodedData))
	if err != nil {
		log.Fatalf("Failed to put object: %v", err)
		return "", err
	}
	signedURL := FullUrlWithDomain(ctx, fmt.Sprintf("%s/%s", common.RESOURCE_AVATAR, objectKey))
	return signedURL, nil
}
