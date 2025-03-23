package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// WechatDecrypt 解密微信加密数据
func WechatDecrypt(sessionKey, encryptedData, iv string) ([]byte, error) {
	// Base64解码
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, errors.New("Base64解码sessionKey失败: " + err.Error())
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.New("Base64解码iv失败: " + err.Error())
	}

	cryptData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.New("Base64解码encryptedData失败: " + err.Error())
	}

	// AES-128-CBC解密
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, errors.New("创建AES解密器失败: " + err.Error())
	}

	// 校验数据长度
	if len(cryptData) < aes.BlockSize {
		return nil, errors.New("加密数据长度不足")
	}

	// 校验初始向量长度
	if len(ivBytes) != aes.BlockSize {
		return nil, errors.New("初始向量长度不正确")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	plainData := make([]byte, len(cryptData))
	mode.CryptBlocks(plainData, cryptData)

	// 去掉PKCS#7填充
	paddingLen := int(plainData[len(plainData)-1])
	if paddingLen < 1 || paddingLen > aes.BlockSize {
		return nil, errors.New("无效的PKCS#7填充")
	}

	return plainData[:len(plainData)-paddingLen], nil
}
