package cacher

import "time"

type CacheConfig struct {
	Path         string   `yaml:"path"` // 路径
	Driver       string   `yaml:"driver"`
	Password     string   `yaml:"password"`
	Prefix       string   `yaml:"prefix"`
	Addr         []string `yaml:"addr"`
	DB           int      `yaml:"db"`
	MinIdleConns int      `yaml:"minIdleConns"`
	PoolSize     int      `yaml:"poolSize"`
}

// 缓存接口
type Cacher interface {
	// 保存数据
	Set(key string, value any, expiration time.Duration) bool
	// 保存数据
	SetNX(key string, value any, expiration time.Duration) bool
	SetXX(key string, value any, expiration time.Duration) bool
	// 累加
	Incr(key string) int64
	// 累加
	IncrBy(key string, val int64) int64
	Decr(key string) int64
	DecrBy(key string, val int64) int64
	// KEY是否存在
	Exists(keys ...string) int64
	// 获取数据 string
	Get(key string) string
	MGet(keys ...string) []any
	// 获取数据 bytes
	GetBytes(key string) []byte
	// 获取数据 int
	GetInt(key string) int
	// 获取数据 int64
	GetInt64(key string) int64
	// 持久化 不过期存储
	GetSet(key string, value any) string
	// 批量获取KEY
	GetPatternKeys(prefix string) []string
	// 批量获取KEY
	GetPatternScan(prefix string) []string
	// 自动加前缀 批量删除KEY
	Delete(keys ...string) bool
	Unlink(keys ...string) bool
	// 无前缀 批量删除KEY
	DeleteKeys(keys ...string) bool
	UnlinkKeys(keys ...string) bool
	// 加锁
	LockStart(key string, args ...int) bool
	// 解锁
	LockEnd(key string)
	// 关闭
	Close()
	HExists(key, field string) bool
	HGet(key, field string) string
	HGetBytes(key, field string) []byte
	HGetInt64(key, field string) int64
	HGetAll(key string) map[string]string
	HKeys(key string) []string
	HSet(key string, values ...any) int64
	HDel(key string, fields ...string) int64
	HLen(key string) int64
	HIncrBy(key, field string, val int64) int64
	FlushDB() bool

	Expire(key string, expiration time.Duration) bool
	PExpire(key string, expiration time.Duration) bool
	ExpireAt(key string, tm time.Time) bool
	PExpireAt(key string, tm time.Time) bool

	GetBit(key string, offset int64) int64
	SetBit(key string, offset int64, val int) int64
	BitCount(key string, start, end int64) int64

	LLen(key string) int64
	LPush(key string, values ...any) int64
	LPop(key string) []byte
	RPush(key string, values ...any) int64
	RPop(key string) []byte
	LRange(key string, start, stop int64) []string
}
