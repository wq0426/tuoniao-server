// @program:     countrybattle
// @file:        base.go
// @author:      ac
// @create:      2024-11-05 17:51
// @description:
package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"app/pkg/config"
)

type Cache struct {
	Ctx     context.Context
	Client  *redis.Client
	Prefix  string
	Expired time.Duration
}

func NewCache(ctx context.Context, prefix string) *Cache {
	return &Cache{
		Ctx:     ctx,
		Client:  config.Rdb,
		Prefix:  prefix,
		Expired: 15 * time.Minute,
	}
}

func (c *Cache) Set(key string, value any) error {
	key = c.Prefix + key
	err := c.Client.Set(c.Ctx, key, value, c.Expired).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Del(key string) error {
	key = c.Prefix + key
	err := c.Client.Del(c.Ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetString(key string) (string, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.Get(c.Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return categoryDeviceStr, nil
}

func (c *Cache) GetInt(key string) (int, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.Get(c.Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	categoryDeviceInt, err := strconv.Atoi(categoryDeviceStr)
	if err != nil {
		return 0, err
	}
	return categoryDeviceInt, nil
}

func (c *Cache) GetBool(key string) (bool, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.Get(c.Ctx, key).Result()
	if err != nil {
		return false, err
	}
	categoryDeviceBool, err := strconv.ParseBool(categoryDeviceStr)
	if err != nil {
		return false, err
	}
	return categoryDeviceBool, nil
}

func (c *Cache) HGetString(key string, field string) (string, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.HGet(c.Ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return categoryDeviceStr, nil
}

func (c *Cache) HGetInt(key string, field string) (int, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.HGet(c.Ctx, key, field).Result()
	if err != nil {
		return 0, err
	}
	categoryDeviceInt, err := strconv.Atoi(categoryDeviceStr)
	if err != nil {
		return 0, err
	}
	return categoryDeviceInt, nil
}

func (c *Cache) HGetBool(key string, field string) (bool, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.HGet(c.Ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	categoryDeviceBool, err := strconv.ParseBool(categoryDeviceStr)
	if err != nil {
		return false, err
	}
	return categoryDeviceBool, nil
}

func (c *Cache) HGetAll(key string) (map[string]string, error) {
	key = c.Prefix + key
	categoryDeviceStr, err := c.Client.HGetAll(c.Ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return categoryDeviceStr, nil
}

func (c *Cache) HSet(key string, values ...any) error {
	key = c.Prefix + key
	err := c.Client.HSet(c.Ctx, key, values).Err()
	if err != nil {
		return err
	}
	err = c.Client.Expire(c.Ctx, key, c.Expired).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) HDel(key string, field string) error {
	key = c.Prefix + key
	err := c.Client.HDel(c.Ctx, key, field).Err()
	if err != nil {
		return err
	}
	return nil
}
