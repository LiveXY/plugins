package plugin

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/livexy/plugin/cacher"

	"github.com/livexy/linq"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

const lockSeconds = 10

type redisCache struct {
	rdb    redis.UniversalClient
	logger *zap.Logger
	prefix string
}

// NewRedisCache 创建一个新的 Redis 缓存实例
// 支持单节点和集群模式，根据 AppConfig 中的配置自动切换
func NewRedisCache(cfg cacher.CacheConfig, logger *zap.Logger) (cacher.Cacher, error) {
	cache := &redisCache{logger: logger, prefix: cfg.Prefix}
	if len(cfg.Addr) == 0 {
		return nil, errors.New("请在config.yaml中配置cache缓存")
	}

	var rdb redis.UniversalClient
	if len(cfg.Addr) == 1 {
		redisOptions := &redis.Options{
			Addr: cfg.Addr[0], DB: cfg.DB,
			MinIdleConns: cfg.MinIdleConns,
			PoolSize:     cfg.PoolSize,
		}
		if len(cfg.Password) > 0 {
			redisOptions.Password = cfg.Password
		}
		rdb = redis.NewClient(redisOptions)
	} else {
		redisOptions := &redis.ClusterOptions{
			Addrs: cfg.Addr, MinIdleConns: cfg.MinIdleConns,
			PoolSize: cfg.PoolSize,
		}
		if len(cfg.Password) > 0 {
			redisOptions.Password = cfg.Password
		}
		rdb = redis.NewClusterClient(redisOptions)
	}

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	cache.rdb = rdb
	return cache, nil
}

// 保存数据
func (cache *redisCache) Set(key string, value any, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	err := cache.rdb.Set(ctx, ckey, value, expiration).Err()
	if err != nil {
		cache.logger.Error("Redis Set：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
		return false
	}
	return true
}

// 保存数据
func (cache *redisCache) SetNX(key string, value any, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	result, err := cache.rdb.SetNX(ctx, ckey, value, expiration).Result()
	if err != nil {
		cache.logger.Error("Redis SetNX：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
		return false
	}
	return result
}
func (cache *redisCache) SetXX(key string, value any, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	result, err := cache.rdb.SetXX(ctx, ckey, value, expiration).Result()
	if err != nil {
		cache.logger.Error("Redis SetXX：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
		return false
	}
	return result
}

// 保存数据
func (cache *redisCache) Incr(key string) int64 {
	ckey := cache.prefix + ":" + key
	result, err := cache.rdb.Incr(ctx, ckey).Result()
	if err != nil {
		cache.logger.Error("Redis Incr：", zap.String("key", key), zap.Error(err))
		return 0
	}
	return result
}
func (cache *redisCache) Decr(key string) int64 {
	ckey := cache.prefix + ":" + key
	result, err := cache.rdb.Decr(ctx, ckey).Result()
	if err != nil {
		cache.logger.Error("Redis Decr：", zap.String("key", key), zap.Error(err))
		return 0
	}
	return result
}

// 保存数据
func (cache *redisCache) IncrBy(key string, val int64) int64 {
	ckey := cache.prefix + ":" + key
	result, err := cache.rdb.IncrBy(ctx, ckey, val).Result()
	if err != nil {
		cache.logger.Error("Redis IncrBy：", zap.String("key", key), zap.Int64("value", val), zap.Error(err))
		return 0
	}
	return result
}
func (cache *redisCache) DecrBy(key string, val int64) int64 {
	ckey := cache.prefix + ":" + key
	result, err := cache.rdb.DecrBy(ctx, ckey, val).Result()
	if err != nil {
		cache.logger.Error("Redis DecrBy：", zap.String("key", key), zap.Int64("value", val), zap.Error(err))
		return 0
	}
	return result
}

// KEY是否存在
func (cache *redisCache) Exists(keys ...string) int64 {
	for i := range keys {
		keys[i] = cache.prefix + ":" + keys[i]
	}
	result, err := cache.rdb.Exists(ctx, keys...).Result()
	if err != nil {
		cache.logger.Error("Redis Exists：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
		return 0
	}
	return result
}

// 获取数据 string
func (cache *redisCache) Get(key string) string {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.Get(ctx, ckey).Result()
	if err != nil {
		return ""
	}
	return val
}
func (cache *redisCache) MGet(ks ...string) []any {
	keys := make([]string, 0, len(ks))
	for _, v := range ks {
		keys = append(keys, cache.prefix+":"+v)
	}
	val, err := cache.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		cache.logger.Error("Redis MGet：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
		return nil
	}
	return val
}

// 获取数据 bytes
func (cache *redisCache) GetBytes(key string) []byte {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.Get(ctx, ckey).Bytes()
	if err != nil {
		return nil
	}
	return val
}

// 获取INT数据
func (cache *redisCache) GetInt(key string) int {
	sval := cache.Get(key)
	ival, err := strconv.Atoi(sval)
	if err != nil {
		ival = 0
	}
	return ival
}

// 获取INT64数据
func (cache *redisCache) GetInt64(key string) int64 {
	sval := cache.Get(key)
	ival, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		ival = 0
	}
	return ival
}

// 持久化 不过期存储
func (cache *redisCache) GetSet(key string, value any) string {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.GetSet(ctx, ckey, value).Result()
	if err != nil {
		cache.logger.Error("Redis GetSet：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
		return ""
	}
	return val
}

// 批量获取KEY
func (cache *redisCache) GetPatternKeys(prefix string) []string {
	key := prefix + "*"
	return cache.rdb.Keys(ctx, cache.prefix+":"+key).Val()
}

// 批量获取KEY
func (cache *redisCache) GetPatternScan(prefix string) []string {
	list := []string{}
	key := prefix + "*"
	var cursor uint64
	for {
		keys, cur, err := cache.rdb.Scan(ctx, cursor, cache.prefix+":"+key, 20).Result()
		if err != nil {
			cache.logger.Error("Redis GetPatternScan：", zap.String("prefix", prefix), zap.Error(err))
			break
		}
		if len(keys) > 0 {
			list = append(list, keys...)
		}
		cursor = cur
		if cur == 0 {
			break
		}
	}
	return linq.Uniq(list)
}

// 自动加前缀 批量删除KEY
func (cache *redisCache) Delete(keys ...string) bool {
	for i, v := range keys {
		keys[i] = cache.prefix + ":" + v
	}
	_, err := cache.rdb.Del(ctx, keys...).Result()
	if err != nil {
		cache.logger.Error("Redis Delete：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
		return false
	}
	return true
}
func (cache *redisCache) Unlink(keys ...string) bool {
	for i, v := range keys {
		keys[i] = cache.prefix + ":" + v
	}
	_, err := cache.rdb.Unlink(ctx, keys...).Result()
	if err != nil {
		cache.logger.Error("Redis Unlink：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
		return false
	}
	return true
}

// 无前缀 批量删除KEY
func (cache *redisCache) DeleteKeys(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	_, err := cache.rdb.Del(ctx, keys...).Result()
	if err != nil {
		cache.logger.Error("Redis DeleteKeys：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
		return false
	}
	return true
}
func (cache *redisCache) UnlinkKeys(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	_, err := cache.rdb.Unlink(ctx, keys...).Result()
	if err != nil {
		cache.logger.Error("Redis UnlinkKeys：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
		return false
	}
	return true
}

// 加锁
func (cache *redisCache) LockStart(key string, args ...int) bool {
	seconds := lockSeconds
	if len(args) > 0 {
		seconds = args[0]
	}
	key = "Lock:" + key
	return !cache.SetNX(key, 1, time.Duration(seconds)*time.Second)
}

// 解锁
func (cache *redisCache) LockEnd(key string) {
	key = "Lock:" + key
	cache.Delete(key)
}

// 关闭释放连接
func (cache *redisCache) Close() {
	if cache.rdb != nil {
		if err := cache.rdb.Close(); err != nil {
			cache.logger.Error("关闭连接失败：", zap.Error(err))
		}
	}
}

// 存在
func (cache *redisCache) HExists(key, field string) bool {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HExists(ctx, ckey, field).Result()
	if err != nil {
		cache.logger.Error("Redis HExists：", zap.String("key", key), zap.String("field", field), zap.Error(err))
		return false
	}
	return val
}

// 获取数据 string
func (cache *redisCache) HGet(key, field string) string {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HGet(ctx, ckey, field).Result()
	if err != nil {
		//cache.logger.Error("Redis HGet：", zap.String("key", key), zap.String("field", field), zap.Error(err))
		return ""
	}
	return val
}

// 获取数据 bytes
func (cache *redisCache) HGetBytes(key, field string) []byte {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HGet(ctx, ckey, field).Bytes()
	if err != nil {
		//cache.logger.Error("Redis HGetBytes：", zap.String("key", key), zap.String("field", field), zap.Error(err))
		return nil
	}
	return val
}

// 获取INT64数据
func (cache *redisCache) HGetInt64(key, field string) int64 {
	sval := cache.HGet(key, field)
	ival, err := strconv.ParseInt(sval, 10, 64)
	if err != nil {
		ival = 0
	}
	return ival
}
func (cache *redisCache) HGetAll(key string) map[string]string {
	all := make(map[string]string)
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HGetAll(ctx, ckey).Result()
	if err != nil {
		return all
	}
	all = val
	return all
}
func (cache *redisCache) HKeys(key string) []string {
	var all []string
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HKeys(ctx, ckey).Result()
	if err != nil {
		cache.logger.Error("Redis HKeys：", zap.String("key", key), zap.Error(err))
		return all
	}
	all = val
	return all
}
func (cache *redisCache) HLen(key string) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HLen(ctx, ckey).Result()
	if err != nil {
		cache.logger.Error("Redis HLen：", zap.String("key", key), zap.Error(err))
		return 0
	}
	return val
}

// 保存
func (cache *redisCache) HSet(key string, values ...any) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HSet(ctx, ckey, values...).Result()
	if err != nil {
		cache.logger.Error("Redis HSet：", zap.String("key", key), zap.Any("value", values), zap.Error(err))
		return 0
	}
	return val
}
func (cache *redisCache) HIncrBy(key, field string, incr int64) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HIncrBy(ctx, ckey, field, incr).Result()
	if err != nil {
		cache.logger.Error("Redis HIncrBy：", zap.String("key", key), zap.String("field", field), zap.Int64("incr", incr), zap.Error(err))
		return 0
	}
	return val
}
func (cache *redisCache) HDel(key string, fields ...string) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.HDel(ctx, ckey, fields...).Result()
	if err != nil {
		cache.logger.Error("Redis HDel：", zap.String("key", key), zap.String("field", strings.Join(fields, ";")), zap.Error(err))
		return 0
	}
	return val
}

func (cache *redisCache) FlushDB() bool {
	_, err := cache.rdb.FlushDB(ctx).Result()
	if err != nil {
		cache.logger.Error("Redis FlushDB：", zap.Error(err))
		return false
	}
	return true
}

func (cache *redisCache) Expire(key string, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.Expire(ctx, ckey, expiration).Result()
	if err != nil {
		cache.logger.Error("Redis Expire", zap.String("key", key), zap.Duration("exp", expiration), zap.Error(err))
		return false
	}
	return val
}
func (cache *redisCache) PExpire(key string, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.PExpire(ctx, ckey, expiration).Result()
	if err != nil {
		cache.logger.Error("Redis PExpire", zap.String("key", key), zap.Duration("exp", expiration), zap.Error(err))
		return false
	}
	return val
}
func (cache *redisCache) ExpireAt(key string, tm time.Time) bool {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.ExpireAt(ctx, ckey, tm).Result()
	if err != nil {
		cache.logger.Error("Redis ExpireAt", zap.String("key", key), zap.Time("tm", tm), zap.Error(err))
		return false
	}
	return val
}
func (cache *redisCache) PExpireAt(key string, tm time.Time) bool {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.PExpireAt(ctx, ckey, tm).Result()
	if err != nil {
		cache.logger.Error("Redis PExpireAt", zap.String("key", key), zap.Time("tm", tm), zap.Error(err))
		return false
	}
	return val
}

func (cache *redisCache) GetBit(key string, offset int64) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.GetBit(ctx, ckey, offset).Result()
	if err != nil {
		cache.logger.Error("Redis GetBit", zap.String("key", key), zap.Int64("offset", offset), zap.Error(err))
		return 0
	}
	return val
}
func (cache *redisCache) SetBit(key string, offset int64, val int) int64 {
	ckey := cache.prefix + ":" + key
	val64, err := cache.rdb.SetBit(ctx, ckey, offset, val).Result()
	if err != nil {
		cache.logger.Error("Redis SetBit", zap.String("key", key), zap.Int64("offset", offset), zap.Int("val", val), zap.Error(err))
		return 0
	}
	return val64
}
func (cache *redisCache) BitCount(key string, start, end int64) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.BitCount(ctx, ckey, &redis.BitCount{Start: start, End: end}).Result()
	if err != nil {
		cache.logger.Error("Redis BitCount", zap.String("key", key), zap.Int64("start", start), zap.Int64("end", end), zap.Int64("val", val), zap.Error(err))
		return 0
	}
	return val
}

func (cache *redisCache) LLen(key string) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.LLen(ctx, ckey).Result()
	if err != nil {
		cache.logger.Error("Redis LLen", zap.String("key", key), zap.Error(err))
		return 0
	}
	return val
}
func (cache *redisCache) LPush(key string, values ...any) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.LPush(ctx, ckey, values...).Result()
	if err != nil {
		cache.logger.Error("Redis LPush", zap.String("key", key), zap.Error(err))
		return 0
	}
	return val
}
func (cache *redisCache) LPop(key string) []byte {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.LPop(ctx, ckey).Bytes()
	if err != nil {
		cache.logger.Error("Redis LPop", zap.String("key", key), zap.Error(err))
		return nil
	}
	return val
}
func (cache *redisCache) RPush(key string, values ...any) int64 {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.RPush(ctx, ckey, values...).Result()
	if err != nil {
		cache.logger.Error("Redis RPush", zap.String("key", key), zap.Error(err))
		return 0
	}
	return val
}
func (cache *redisCache) RPop(key string) []byte {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.RPop(ctx, ckey).Bytes()
	if err != nil {
		cache.logger.Error("Redis RPop", zap.String("key", key), zap.Error(err))
		return nil
	}
	return val
}
func (cache *redisCache) LRange(key string, start, stop int64) []string {
	ckey := cache.prefix + ":" + key
	val, err := cache.rdb.LRange(ctx, ckey, start, stop).Result()
	if err != nil {
		cache.logger.Error("Redis LRange", zap.String("key", key), zap.Int64("start", start), zap.Int64("stop", stop), zap.Error(err))
		return nil
	}
	return val
}
