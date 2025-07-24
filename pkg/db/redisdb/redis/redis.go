package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"liveJob/pkg/common/config"
)

var (
	rdb *Db
	Nil = redis.Nil
)

type Db struct {
	wPool redis.UniversalClient
}

func InitRedis() error {
	rwPool, err := initRedis(strings.Join(config.Config.Redis.Address, ","), config.Config.Redis.Password, "", 50)
	if err != nil {
		return err
	}

	rdb = &Db{
		wPool: rwPool,
	}
	return nil
}

func InitRedisByConfig(address []string, password string) error {
	rwPool, err := initRedis(strings.Join(address, ","), password, "", 50)
	if err != nil {
		return err
	}

	rdb = &Db{
		wPool: rwPool,
	}
	return nil
}

func initRedis(host, auth, master string, poolSize int) (redis.UniversalClient, error) {
	options := &redis.UniversalOptions{
		Addrs:           strings.Split(host, ","), // redis地址
		MaxRedirects:    0,                        // 放弃前最大重试次数,默认是不重试失败的命令,默认是3次
		ReadOnly:        false,                    // 在从库上打开只读命令
		RouteByLatency:  false,                    // 允许将只读命令路由到最近的主节点或从节点,自动启用只读
		RouteRandomly:   false,                    // 允许将只读命令路由到随机主节点或从节点。 它自动启用只读。
		Password:        auth,
		MaxRetries:      2,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    20 * time.Second,
		PoolSize:        poolSize,
		PoolTimeout:     30 * time.Second,
	}
	// 哨兵模式
	if len(master) > 0 {
		options.SentinelPassword = auth
		options.MasterName = master
	}
	redisPool := redis.NewUniversalClient(options)
	_, err := redisPool.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return redisPool, nil
}

func GetKey(key string) (string, error) {
	value, err := rdb.wPool.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return value, nil
}

func Eval(script string, keys []string, args ...interface{}) error {
	return rdb.wPool.Eval(context.Background(), script, keys, args...).Err()
}

func GetKeyBytes(key string) ([]byte, error) {
	return rdb.wPool.Get(context.Background(), key).Bytes()
}

// SetNotExpireKV 设置不过期的 key
func SetNotExpireKV(key, value string) error {
	err := rdb.wPool.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetNotExpireKVInterface 设置不过期的 key
func SetNotExpireKVInterface(key string, value interface{}) error {
	err := rdb.wPool.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetExpireKV 设置过期的 key
func SetExpireKV(key, value string, expire time.Duration) error {
	err := rdb.wPool.Set(context.Background(), key, value, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetTTL(key string) time.Duration {
	return rdb.wPool.TTL(context.Background(), key).Val()
}

// SetExpireKey 设置 key 过期
func SetExpireKey(key string, expire time.Duration) error {
	err := rdb.wPool.Expire(context.Background(), key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

// Set 设置 key, value 以及过期时间
func Set(key string, value string, expire time.Duration) error {
	_, err := rdb.wPool.Set(context.Background(), key, value, expire).Result()
	if err != nil {
		return err
	}
	return nil
}

// SetNX 设置 key, value 以及过期时间
func SetNX(key string, value string, expire time.Duration) (bool, error) {
	flag, err := rdb.wPool.SetNX(context.Background(), key, value, expire).Result()
	if err != nil {
		return false, err
	}
	return flag, nil
}

// DelKey 删除 redis 的key
func DelKey(key string) error {
	return rdb.wPool.Del(context.Background(), key).Err()
}

// KeyExist 判断某一个key 是否存在
func KeyExist(keys string) (bool, error) {

	count, err := rdb.wPool.Exists(context.Background(), keys).Result()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// HSet 设置 hash
func HSet(key, field string, value interface{}) error {
	return rdb.wPool.HSet(context.Background(), key, field, value).Err()
}

// HMSet 批量存储 hash
func HMSet(key string, fields map[string]interface{}) error {
	if len(fields) < 1 {
		return nil
	}

	err := rdb.wPool.HMSet(context.Background(), key, fields).Err()
	if err != nil {
		return err
	}

	return nil
}

// HGet 获取单个 hash
func HGet(key, field string) (string, error) {
	return rdb.wPool.HGet(context.Background(), key, field).Result()
}

func HKeys(key string) ([]string, error) {
	return rdb.wPool.HKeys(context.Background(), key).Result()
}

// HMGet 批量获取 hash
func HMGet(key string, fields ...string) ([]interface{}, error) {
	res, err := rdb.wPool.HMGet(context.Background(), key, fields...).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HGetAll 获取 hash 全部值
func HGetAll(key string) (map[string]string, error) {
	res, err := rdb.wPool.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HScan 获取 hash 键值树
func HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rdb.wPool.HScan(context.Background(), key, cursor, match, count).Result()
}

func SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rdb.wPool.SScan(context.Background(), key, cursor, match, count).Result()
}

func HLen(key string) (int, error) {
	res, err := rdb.wPool.HLen(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

// HDel 删除 hash key
func HDel(key string, fields ...string) error {
	err := rdb.wPool.HDel(context.Background(), key, fields...).Err()
	if err != nil {
		return err
	}

	return nil
}

// HIncr 自增1
func HIncr(key, field string) {
	rdb.wPool.HIncrBy(context.Background(), key, field, 1)
}

// HIncrBy 自增
func HIncrBy(key, field string, incr int64) (int64, error) {
	return rdb.wPool.HIncrBy(context.Background(), key, field, incr).Result()
}

// RPush 在名称为key的list尾添加一个值为value的元素
func RPush(key string, values ...interface{}) error {
	return rdb.wPool.RPush(context.Background(), key, values...).Err()
}

// LPush 在名称为key的list头添加一个值为value的 元素
func LPush(key string, values ...interface{}) error {
	return rdb.wPool.LPush(context.Background(), key, values...).Err()
}

// Publish 在名称为key的list头添加一个值为value的 元素
func Publish(channel string, values interface{}) error {
	return rdb.wPool.Publish(context.Background(), channel, values).Err()
}

// LLen 返回名称为key的list的长度
func LLen(key string) (int64, error) {
	return rdb.wPool.LLen(context.Background(), key).Result()
}

// LRange 返回名称为key的list中start至end之间的元素, start为0, end为-1 则是获取所有 list key
func LRange(key string, start, end int64) ([]string, error) {
	return rdb.wPool.LRange(context.Background(), key, start, end).Result()
}

// LSet 给名称为key的list中index位置的元素赋值
func LSet(key string, index int64, value interface{}) error {
	return rdb.wPool.LSet(context.Background(), key, index, value).Err()
}

// LRem 删除count个key的list中值为value的元素
func LRem(key string, count int64, value interface{}) error {
	return rdb.wPool.LRem(context.Background(), key, count, value).Err()
}

// ZCount  有序集合中 min-max中的成员数量
func ZCount(key, min, max string) (int64, error) {
	count, err := rdb.wPool.ZCount(context.Background(), key, min, max).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ZAdd 获取中元素的数量
func ZAdd(key string, members ...redis.Z) (int64, error) {
	count, err := rdb.wPool.ZAdd(context.Background(), key, members...).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ZCARD 获取中元素的数量
func ZCARD(key string) (int64, error) {
	count, err := rdb.wPool.ZCard(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ZRange 通过索引区间返回有序集合成指定区间内的成员
func ZRange(key string, start, stop int64) ([]string, error) {
	arr, err := rdb.wPool.ZRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return arr, nil
}

// ZRangeByScore 通过索引区间返回有序集合成指定区间内的成员
func ZRangeByScore(key string, min, max string) ([]string, error) {
	opt := redis.ZRangeBy{
		Min: min,
		Max: max,
	}
	arr, err := rdb.wPool.ZRangeByScore(context.Background(), key, &opt).Result()
	if err != nil {
		return []string{}, err
	}

	return arr, nil
}

// ZRevRangeByScore 通过索引区间返回有序集合成指定区间内的成员
func ZRevRangeByScore(key string, min, max string) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: 0,
		Count:  1,
	}
	arr, err := rdb.wPool.ZRevRangeByScore(context.Background(), key, &opt).Result()
	if err != nil {
		return []string{}, err
	}

	return arr, nil
}

// ZRangeByScorePageInfo 通过索引区间分页返回有序集合成指定区间内的成员
func ZRangeByScorePageInfo(key string, min, max string, pageNo, pageSize int64) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: (pageNo - 1) * pageSize,
		Count:  pageSize,
	}
	arr, err := rdb.wPool.ZRangeByScore(context.Background(), key, &opt).Result()
	if err != nil {
		return []string{}, err
	}
	return arr, nil
}

func ZRem(key string, members ...string) error {
	return rdb.wPool.ZRem(context.Background(), key, members).Err()
}

// ZRangeWithScores 返回所有有序集合
func ZRangeWithScores(key string) ([]redis.Z, error) {
	return rdb.wPool.ZRangeWithScores(context.Background(), key, 0, -1).Result()
}

func ZScore(key string, member string) (float64, error) {
	return rdb.wPool.ZScore(context.Background(), key, member).Result()
}

func HGetBytesByField(key, filed string) ([]byte, error) {
	return rdb.wPool.HGet(context.Background(), key, filed).Bytes()
}

func SIsMember(key, field string) (bool, error) {
	return rdb.wPool.SIsMember(context.Background(), key, field).Result()
}
func Incr(key string) {
	rdb.wPool.Incr(context.Background(), key)
}

func IncrBy(key string, value int64) (int64, error) {
	return rdb.wPool.IncrBy(context.Background(), key, value).Result()
}

func IncrWithResult(key string) (int64, error) {
	return rdb.wPool.Incr(context.Background(), key).Result()
}

func DecrWithResult(key string) (int64, error) {
	return rdb.wPool.Decr(context.Background(), key).Result()
}

func SMembers(key string) ([]string, error) {
	return rdb.wPool.SMembers(context.Background(), key).Result()
}

func SAdd(key string, members ...interface{}) (int64, error) {
	return rdb.wPool.SAdd(context.Background(), key, members...).Result()
}

func SRem(key string, members ...interface{}) (int64, error) {
	return rdb.wPool.SRem(context.Background(), key, members...).Result()
}

func LIndex(key string, index int64) (string, error) {
	return rdb.wPool.LIndex(context.Background(), key, index).Result()
}

func SPop(key string) (string, error) {
	return rdb.wPool.SPop(context.Background(), key).Result()
}

// SetSscan 集合读取
func SetSscan(key string, match string, perCount int64) ([]string, error) {
	var (
		cursor = uint64(0)
		data   []string
	)
	for {
		keys, retCursor, err := rdb.wPool.SScan(context.Background(), key, cursor, match, perCount).Result()
		if err != nil {
			return data, err
		}
		if len(keys) == 0 {
			break
		}
		data = append(data, keys...)
		if retCursor == 0 {
			break
		}
		cursor = retCursor
	}
	return data, nil
}

// 获取键值,如不存在 则获取func 存入到键中
func GetOrSet(key string, f func() (interface{}, error), expire time.Duration) ([]byte, error) {
	result, err := rdb.wPool.Get(context.Background(), key).Bytes()
	if err != nil || len(result) == 0 {
		data, err := f()
		if err == nil {
			var value []byte
			value, err = json.Marshal(data)
			if err != nil {
				return nil, err
			}
			err = rdb.wPool.Set(context.Background(), key, value, expire).Err()
			if err != nil {
				return nil, err
			}
			return value, nil
		}
		return nil, err
	}
	return result, nil
}
func GetRedisPool() redis.UniversalClient {
	return rdb.wPool
}
