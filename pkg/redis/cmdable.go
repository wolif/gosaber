package redis

import (
	"time"
)

func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(key, value, expiration).Err()
}

func Get(key string) (string, error) {
	stringCmd := Client.Get(key)
	return stringCmd.Val(), stringCmd.Err()
}

func MGet(keys []string) ([]interface{}, error) {
	sliceCmd := Client.MGet(keys...)
	return sliceCmd.Val(), sliceCmd.Err()
}

func HSet(key, field string, value interface{}) error {
	return Client.HSet(key, field, value).Err()
}

func HIncrBy(key, field string, incr int64) (int64, error) {
	intCmd := Client.HIncrBy(key, field, incr)
	return intCmd.Val(), intCmd.Err()
}

func HGet(key, field string) (string, error) {
	stringCmd := Client.HGet(key, field)
	return stringCmd.Val(), stringCmd.Err()
}

func HGetAll(key string) (map[string]string, error) {
	stringStringMapCmd := Client.HGetAll(key)
	return stringStringMapCmd.Val(), stringStringMapCmd.Err()
}

func HKeys(key string) []string {
	return Client.HKeys(key).Val()
}

func HMGet(key string, fields ...string) []interface{} {
	return Client.HMGet(key, fields...).Val()
}

func HDel(key string, field ...string) error {
	return Client.HDel(key, field...).Err()
}

func Expire(key string, t time.Duration) error {
	return Client.Expire(key, t).Err()
}

func HSetTTL(key, field string, value interface{}, t time.Duration) error {
	if err := Client.HSet(key, field, value).Err(); err != nil {
		return err
	}
	return Client.Expire(key, t).Err()
}

func Keys(key string) []string {
	return Client.Keys(key).Val()
}

func Del(keys ...string) error {
	return Client.Del(keys...).Err()
}
