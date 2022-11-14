package plugin

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/livexy/plugins/plugin/cacher"

	"github.com/livexy/linq"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var ctx = context.Background()

const lockSeconds = 10

type redisCache struct {
	client  *redis.Client
	cluster *redis.ClusterClient
	logger *zap.Logger
	prefix  string
}

// REDIS实例
func NewRedisCache(cfg cacher.CacheConfig, logger *zap.Logger) (cacher.Cacher, error) {
	cache := &redisCache{ logger: logger, prefix: cfg.Prefix }
	if len(cfg.Addr) == 0 {
		return nil, errors.New("请在config.yaml中配置cache缓存")
	}
	if len(cfg.Addr) == 1 {
		redisOptions := &redis.Options{
			Addr: cfg.Addr[0], DB: cfg.DB,
			MinIdleConns: cfg.MinIdleConns,
			PoolSize:     cfg.PoolSize,
		}
		if len(cfg.Password) > 0 {
			redisOptions.Password = cfg.Password
		}
		cache.client = redis.NewClient(redisOptions)
		_, err := cache.client.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}
	} else {
		redisOptions := &redis.ClusterOptions{
			Addrs: cfg.Addr, MinIdleConns: cfg.MinIdleConns,
			PoolSize: cfg.PoolSize,
		}
		if len(cfg.Password) > 0 {
			redisOptions.Password = cfg.Password
		}
		cache.cluster = redis.NewClusterClient(redisOptions)
		_, err := cache.cluster.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}
	}
	return cache, nil
}

// 保存数据
func (cache *redisCache) Set(key string, value any, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		err := cache.cluster.Set(ctx, ckey, value, expiration).Err()
		if err != nil {
			cache.logger.Error("RedisCluster Set：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return false
		}
	} else {
		err := cache.client.Set(ctx, ckey, value, expiration).Err()
		if err != nil {
			cache.logger.Error("Redis Set：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return false
		}
	}
	return true
}

// 保存数据
func (cache *redisCache) SetNX(key string, value any, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		result, err := cache.cluster.SetNX(ctx, ckey, value, expiration).Result()
		if err != nil {
			cache.logger.Error("RedisCluster SetNX：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return false
		}
		return result
	} else {
		result, err := cache.client.SetNX(ctx, ckey, value, expiration).Result()
		if err != nil {
			cache.logger.Error("Redis SetNX：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return false
		}
		return result
	}
}
func (cache *redisCache) SetXX(key string, value any, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		result, err := cache.cluster.SetXX(ctx, ckey, value, expiration).Result()
		if err != nil {
			cache.logger.Error("RedisCluster SetXX：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return false
		}
		return result
	} else {
		result, err := cache.client.SetXX(ctx, ckey, value, expiration).Result()
		if err != nil {
			cache.logger.Error("Redis SetXX：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return false
		}
		return result
	}
}

// 保存数据
func (cache *redisCache) Incr(key string) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		result, err := cache.cluster.Incr(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("RedisCluster Incr：", zap.String("key", key), zap.Error(err))
			return 0
		}
		return result
	} else {
		result, err := cache.client.Incr(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("Redis Incr：", zap.String("key", key), zap.Error(err))
			return 0
		}
		return result
	}
}
func (cache *redisCache) Decr(key string) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		result, err := cache.cluster.Decr(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("RedisCluster Decr：", zap.String("key", key), zap.Error(err))
			return 0
		}
		return result
	} else {
		result, err := cache.client.Decr(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("Redis Decr：", zap.String("key", key), zap.Error(err))
			return 0
		}
		return result
	}
}

// 保存数据
func (cache *redisCache) IncrBy(key string, val int64) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		result, err := cache.cluster.IncrBy(ctx, ckey, val).Result()
		if err != nil {
			cache.logger.Error("RedisCluster IncrBy：", zap.String("key", key), zap.Int64("value", val), zap.Error(err))
			return 0
		}
		return result
	} else {
		result, err := cache.client.IncrBy(ctx, ckey, val).Result()
		if err != nil {
			cache.logger.Error("Redis IncrBy：", zap.String("key", key), zap.Int64("value", val), zap.Error(err))
			return 0
		}
		return result
	}
}
func (cache *redisCache) DecrBy(key string, val int64) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		result, err := cache.cluster.DecrBy(ctx, ckey, val).Result()
		if err != nil {
			cache.logger.Error("RedisCluster DecrBy：", zap.String("key", key), zap.Int64("value", val), zap.Error(err))
			return 0
		}
		return result
	} else {
		result, err := cache.client.DecrBy(ctx, ckey, val).Result()
		if err != nil {
			cache.logger.Error("Redis DecrBy：", zap.String("key", key), zap.Int64("value", val), zap.Error(err))
			return 0
		}
		return result
	}
}

// KEY是否存在
func (cache *redisCache) Exists(keys ...string) int64 {
	for i := range keys {
		keys[i] = cache.prefix + ":" + keys[i]
	}
	if cache.cluster != nil {
		result, err := cache.cluster.Exists(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster Exists：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return 0
		}
		return result
	} else {
		result, err := cache.client.Exists(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("Redis Exists：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return 0
		}
		return result
	}
}

// 获取数据 string
func (cache *redisCache) Get(key string) string {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.Get(ctx, ckey).Result()
		if err != nil {
			return ""
		}
		return val
	} else {
		val, err := cache.client.Get(ctx, ckey).Result()
		if err != nil {
			return ""
		}
		return val
	}
}
func (cache *redisCache) MGet(ks ...string) []any {
	keys := make([]string, 0, len(ks))
	for _, v := range ks {
		keys = append(keys, cache.prefix+":"+v)
	}
	if cache.cluster != nil {
		val, err := cache.cluster.MGet(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster MGet：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return nil
		}
		return val
	} else {
		val, err := cache.client.MGet(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("Redis MGet：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return nil
		}
		return val
	}
}

// 获取数据 bytes
func (cache *redisCache) GetBytes(key string) []byte {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.Get(ctx, ckey).Bytes()
		if err != nil {
			return nil
		}
		return val
	} else {
		val, err := cache.client.Get(ctx, ckey).Bytes()
		if err != nil {
			return nil
		}
		return val
	}
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
	if cache.cluster != nil {
		val, err := cache.cluster.GetSet(ctx, ckey, value).Result()
		if err != nil {
			cache.logger.Error("RedisCluster GetSet：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return ""
		}
		return val
	} else {
		val, err := cache.client.GetSet(ctx, ckey, value).Result()
		if err != nil {
			cache.logger.Error("Redis GetSet：", zap.String("key", key), zap.Any("value", value), zap.Error(err))
			return ""
		}
		return val
	}
}

// 批量获取KEY
func (cache *redisCache) GetPatternKeys(prefix string) []string {
	key := prefix + "*"
	if cache.cluster != nil {
		return cache.cluster.Keys(ctx, cache.prefix+":"+key).Val()
	} else {
		return cache.client.Keys(ctx, cache.prefix+":"+key).Val()
	}
}

// 批量获取KEY
func (cache *redisCache) GetPatternScan(prefix string) []string {
	list := []string{}
	key := prefix + "*"
	var cursor uint64
	if cache.cluster != nil {
		for {
			keys, cur, err := cache.cluster.Scan(ctx, cursor, cache.prefix+":"+key, 20).Result()
			if err != nil {
				cache.logger.Error("RedisCluster GetPatternScan：", zap.String("prefix", prefix), zap.Error(err))
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
	} else {
		for {
			keys, cur, err := cache.client.Scan(ctx, cursor, cache.prefix+":"+key, 20).Result()
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
	}
	return linq.Uniq(list)
}

// 自动加前缀 批量删除KEY
func (cache *redisCache) Delete(keys ...string) bool {
	for i, v := range keys {
		keys[i] = cache.prefix + ":" + v
	}
	if cache.cluster != nil {
		_, err := cache.cluster.Del(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster Delete：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	} else {
		_, err := cache.client.Del(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("Redis Delete：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	}
	return true
}
func (cache *redisCache) Unlink(keys ...string) bool {
	for i, v := range keys {
		keys[i] = cache.prefix + ":" + v
	}
	if cache.cluster != nil {
		_, err := cache.cluster.Unlink(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster Unlink：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	} else {
		_, err := cache.client.Unlink(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("Redis Unlink：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	}
	return true
}

// 无前缀 批量删除KEY
func (cache *redisCache) DeleteKeys(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	if cache.cluster != nil {
		_, err := cache.cluster.Del(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster DeleteKeys：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	} else {
		_, err := cache.client.Del(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("Redis DeleteKeys：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	}
	return true
}
func (cache *redisCache) UnlinkKeys(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	if cache.cluster != nil {
		_, err := cache.cluster.Unlink(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster UnlinkKeys：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
	} else {
		_, err := cache.client.Unlink(ctx, keys...).Result()
		if err != nil {
			cache.logger.Error("Redis UnlinkKeys：", zap.String("key", strings.Join(keys, ";")), zap.Error(err))
			return false
		}
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
	if cache.client != nil {
		cache.client.Close()
	}
	if cache.cluster != nil {
		cache.cluster.Close()
	}
}

// 存在
func (cache *redisCache) HExists(key, field string) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HExists(ctx, ckey, field).Result()
		if err != nil {
			cache.logger.Error("RedisCluster HExists：", zap.String("key", key), zap.String("field", field), zap.Error(err))
			return false
		}
		return val
	} else {
		val, err := cache.client.HExists(ctx, ckey, field).Result()
		if err != nil {
			cache.logger.Error("Redis HExists：", zap.String("key", key), zap.String("field", field), zap.Error(err))
			return false
		}
		return val
	}
}

// 获取数据 string
func (cache *redisCache) HGet(key, field string) string {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HGet(ctx, ckey, field).Result()
		if err != nil {
			//cache.logger.Error("RedisCluster HGet：", zap.String("key", key), zap.String("field", field), zap.Error(err))
			return ""
		}
		return val
	} else {
		val, err := cache.client.HGet(ctx, ckey, field).Result()
		if err != nil {
			//cache.logger.Error("Redis HGet：", zap.String("key", key), zap.String("field", field), zap.Error(err))
			return ""
		}
		return val
	}
}

// 获取数据 bytes
func (cache *redisCache) HGetBytes(key, field string) []byte {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HGet(ctx, ckey, field).Bytes()
		if err != nil {
			//cache.logger.Error("RedisCluster HGetBytes：", zap.String("key", key), zap.String("field", field), zap.Error(err))
			return nil
		}
		return val
	} else {
		val, err := cache.client.HGet(ctx, ckey, field).Bytes()
		if err != nil {
			//cache.logger.Error("Redis HGetBytes：", zap.String("key", key), zap.String("field", field), zap.Error(err))
			return nil
		}
		return val
	}
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
	if cache.cluster != nil {
		val, err := cache.cluster.HGetAll(ctx, ckey).Result()
		if err != nil {
			return all
		}
		all = val
		return all
	} else {
		val, err := cache.client.HGetAll(ctx, ckey).Result()
		if err != nil {
			return all
		}
		all = val
		return all
	}
}
func (cache *redisCache) HKeys(key string) []string {
	var all []string
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HKeys(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("RedisCluster HKeys：", zap.String("key", key), zap.Error(err))
			return all
		}
		all = val
		return all
	} else {
		val, err := cache.client.HKeys(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("Redis HKeys：", zap.String("key", key), zap.Error(err))
			return all
		}
		all = val
		return all
	}
}
func (cache *redisCache) HLen(key string) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HLen(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("RedisCluster HLen：", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.HLen(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("Redis HLen：", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	}
}

// 保存
func (cache *redisCache) HSet(key string, values ...any) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HSet(ctx, ckey, values...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster HSet：", zap.String("key", key), zap.Any("value", values), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.HSet(ctx, ckey, values...).Result()
		if err != nil {
			cache.logger.Error("Redis HSet：", zap.String("key", key), zap.Any("value", values), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) HIncrBy(key, field string, incr int64) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HIncrBy(ctx, ckey, field, incr).Result()
		if err != nil {
			cache.logger.Error("RedisCluster HIncrBy：", zap.String("key", key), zap.String("field", field), zap.Int64("incr", incr), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.HIncrBy(ctx, ckey, field, incr).Result()
		if err != nil {
			cache.logger.Error("Redis HIncrBy：", zap.String("key", key), zap.String("field", field), zap.Int64("incr", incr), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) HDel(key string, fields ...string) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.HDel(ctx, ckey, fields...).Result()
		if err != nil {
			cache.logger.Error("RedisCluster HDel：", zap.String("key", key), zap.String("field", strings.Join(fields, ";")), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.HDel(ctx, ckey, fields...).Result()
		if err != nil {
			cache.logger.Error("Redis HDel：", zap.String("key", key), zap.String("field", strings.Join(fields, ";")), zap.Error(err))
			return 0
		}
		return val
	}
}

func (cache *redisCache) FlushDB() bool {
	if cache.cluster != nil {
		_, err := cache.cluster.FlushDB(ctx).Result()
		if err != nil {
			cache.logger.Error("RedisCluster FlushDB：", zap.Error(err))
			return false
		}
	} else {
		_, err := cache.client.FlushDB(ctx).Result()
		if err != nil {
			cache.logger.Error("Redis FlushDB：", zap.Error(err))
			return false
		}
	}
	return true
}

func (cache *redisCache) Expire(key string, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.Expire(ctx, ckey, expiration).Result()
		if err != nil {
			cache.logger.Error("RedisCluster Expire", zap.String("key", key), zap.Duration("exp", expiration), zap.Error(err))
			return false
		}
		return val
	} else {
		val, err := cache.client.Expire(ctx, ckey, expiration).Result()
		if err != nil {
			cache.logger.Error("Redis Expire", zap.String("key", key), zap.Duration("exp", expiration), zap.Error(err))
			return false
		}
		return val
	}
}
func (cache *redisCache) PExpire(key string, expiration time.Duration) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.PExpire(ctx, ckey, expiration).Result()
		if err != nil {
			cache.logger.Error("RedisCluster PExpire", zap.String("key", key), zap.Duration("exp", expiration), zap.Error(err))
			return false
		}
		return val
	} else {
		val, err := cache.client.PExpire(ctx, ckey, expiration).Result()
		if err != nil {
			cache.logger.Error("Redis PExpire", zap.String("key", key), zap.Duration("exp", expiration), zap.Error(err))
			return false
		}
		return val
	}
}
func (cache *redisCache) ExpireAt(key string, tm time.Time) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.ExpireAt(ctx, ckey, tm).Result()
		if err != nil {
			cache.logger.Error("RedisCluster ExpireAt", zap.String("key", key), zap.Time("tm", tm), zap.Error(err))
			return false
		}
		return val
	} else {
		val, err := cache.client.ExpireAt(ctx, ckey, tm).Result()
		if err != nil {
			cache.logger.Error("Redis ExpireAt", zap.String("key", key), zap.Time("tm", tm), zap.Error(err))
			return false
		}
		return val
	}
}
func (cache *redisCache) PExpireAt(key string, tm time.Time) bool {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.PExpireAt(ctx, ckey, tm).Result()
		if err != nil {
			cache.logger.Error("RedisCluster PExpireAt", zap.String("key", key), zap.Time("tm", tm), zap.Error(err))
			return false
		}
		return val
	} else {
		val, err := cache.client.PExpireAt(ctx, ckey, tm).Result()
		if err != nil {
			cache.logger.Error("Redis PExpireAt", zap.String("key", key), zap.Time("tm", tm), zap.Error(err))
			return false
		}
		return val
	}
}

func (cache *redisCache) GetBit(key string, offset int64) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.GetBit(ctx, ckey, offset).Result()
		if err != nil {
			cache.logger.Error("RedisCluster GetBit", zap.String("key", key), zap.Int64("offset", offset), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.GetBit(ctx, ckey, offset).Result()
		if err != nil {
			cache.logger.Error("Redis GetBit", zap.String("key", key), zap.Int64("offset", offset), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) SetBit(key string, offset int64, val int) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.SetBit(ctx, ckey, offset, val).Result()
		if err != nil {
			cache.logger.Error("RedisCluster SetBit", zap.String("key", key), zap.Int64("offset", offset), zap.Int64("val", val), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.SetBit(ctx, ckey, offset, val).Result()
		if err != nil {
			cache.logger.Error("Redis SetBit", zap.String("key", key), zap.Int64("offset", offset), zap.Int64("val", val), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) BitCount(key string, start, end int64) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.BitCount(ctx, ckey, &redis.BitCount{Start: start, End: end}).Result()
		if err != nil {
			cache.logger.Error("RedisCluster BitCount", zap.String("key", key), zap.Int64("start", start), zap.Int64("end", end), zap.Int64("val", val), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.BitCount(ctx, ckey, &redis.BitCount{Start: start, End: end}).Result()
		if err != nil {
			cache.logger.Error("Redis BitCount", zap.String("key", key), zap.Int64("start", start), zap.Int64("end", end), zap.Int64("val", val), zap.Error(err))
			return 0
		}
		return val
	}
}

func (cache *redisCache) LLen(key string) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.LLen(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("RedisCluster LLen", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.LLen(ctx, ckey).Result()
		if err != nil {
			cache.logger.Error("Redis LLen", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) LPush(key string, values ...any) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.LPush(ctx, ckey, values).Result()
		if err != nil {
			cache.logger.Error("RedisCluster LPush", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.LPush(ctx, ckey, values).Result()
		if err != nil {
			cache.logger.Error("Redis LPush", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) LPop(key string) []byte {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.LPop(ctx, ckey).Bytes()
		if err != nil {
			cache.logger.Error("RedisCluster LPop", zap.String("key", key), zap.Error(err))
			return nil
		}
		return val
	} else {
		val, err := cache.client.LPop(ctx, ckey).Bytes()
		if err != nil {
			cache.logger.Error("Redis LPop", zap.String("key", key), zap.Error(err))
			return nil
		}
		return val
	}
}
func (cache *redisCache) RPush(key string, values ...any) int64 {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.RPush(ctx, ckey, values).Result()
		if err != nil {
			cache.logger.Error("RedisCluster RPush", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	} else {
		val, err := cache.client.RPush(ctx, ckey, values).Result()
		if err != nil {
			cache.logger.Error("Redis RPush", zap.String("key", key), zap.Error(err))
			return 0
		}
		return val
	}
}
func (cache *redisCache) RPop(key string) []byte {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.RPop(ctx, ckey).Bytes()
		if err != nil {
			cache.logger.Error("RedisCluster RPop", zap.String("key", key), zap.Error(err))
			return nil
		}
		return val
	} else {
		val, err := cache.client.RPop(ctx, ckey).Bytes()
		if err != nil {
			cache.logger.Error("Redis RPop", zap.String("key", key), zap.Error(err))
			return nil
		}
		return val
	}
}
func (cache *redisCache) LRange(key string, start, stop int64) []string {
	ckey := cache.prefix + ":" + key
	if cache.cluster != nil {
		val, err := cache.cluster.LRange(ctx, ckey, start, stop).Result()
		if err != nil {
			cache.logger.Error("RedisCluster LRange", zap.String("key", key), zap.Int64("start", start), zap.Int64("stop", stop), zap.Error(err))
			return nil
		}
		return val
	} else {
		val, err := cache.client.LRange(ctx, ckey, start, stop).Result()
		if err != nil {
			cache.logger.Error("Redis LRange", zap.String("key", key), zap.Int64("start", start), zap.Int64("stop", stop), zap.Error(err))
			return nil
		}
		return val
	}
}
