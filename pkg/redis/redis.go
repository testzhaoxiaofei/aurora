package redis

import (
	"context"
	"fmt"
	redisV8 "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

type redis struct {
	Client *redisV8.Client
}

var Redis redis
var ctx = context.Background()

func Setup() {
	rDb := redisV8.NewClient(&redisV8.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.Host"), viper.GetString("redis.Port")),
		Password: viper.GetString("redis.Password"),
		DB:       viper.GetInt("redis.Db"),
	})
	_, err := rDb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	Redis = redis{Client: rDb}
}

func (r redis) Exists(key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, err
	}

	return true, nil
}

// HGetAll  获取hash所有字段
func (r redis) HGetAll(key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// Set 设置key
func (r redis) Set(key string, value interface{}, expire time.Duration) (string, error) {
	return r.Client.Set(ctx, key, value, expire*time.Second).Result()
}

func (r redis) SetIncr(key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// TTL 获取
func (r redis) TTL(key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}
func (r redis) Expire(key string, expire time.Duration) (interface{}, error) {
	return r.Client.Expire(ctx, key, expire*time.Second).Result()
}

func (r redis) Get(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r redis) Delete(key string) (int64, error) {
	return r.Client.Del(ctx, key).Result()
}

func (r redis) HINCRBY(key string, field string, value int64) (int64, error) {
	return r.Client.HIncrBy(ctx, key, field, value).Result()
}

func (r redis) HIncrByFloat(key string, field string, value float64) (float64, error) {
	return r.Client.HIncrByFloat(ctx, key, field, value).Result()
}

func (r redis) HINCRBY1(key string, field string, value int64) (int64, error) {
	return r.Client.HIncrBy(ctx, key, field, value).Result()
}

func (r redis) HGET(key string, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

func (r redis) HSET(key string, field string, value interface{}) (int64, error) {
	return r.Client.HSet(ctx, key, field, value).Result()
}
func (r redis) HExists(key string, field string) (bool, error) {
	return r.Client.HExists(ctx, key, field).Result()
}

func (r redis) HDEL(key string, field string) (int64, error) {
	return r.Client.HDel(ctx, key, field).Result()
}

func (r redis) LLEN(key string) (int64, error) {
	return r.Client.LLen(ctx, key).Result()
}

// BLPop 弹出列表中最后一个元素
func (r redis) BRPop(key string) ([]string, error) {
	return r.Client.BRPop(ctx, 2, key).Result()
}

func (r redis) LRange(key string, start int64, end int64) ([]string, error) {
	return r.Client.LRange(ctx, key, start, end).Result()
}

func (r redis) HMGET(key string, k ...string) ([]interface{}, error) {
	return r.Client.HMGet(ctx, key, k...).Result()
}

func (r redis) LPUSH(key string, k string) (interface{}, error) {
	return r.Client.LPush(ctx, key, k).Result()
}
func (r redis) RPush(key string, k string) (interface{}, error) {
	return r.Client.RPush(ctx, key, k).Result()
}
func (r redis) SAdd(key string, k string) (interface{}, error) {
	return r.Client.SAdd(ctx, key, k).Result()
}

func (r redis) SAdds(key string, k interface{}) (interface{}, error) {
	return r.Client.SAdd(ctx, key, k).Result()
}
func (r redis) LushS(key string, k interface{}) (interface{}, error) {
	return r.Client.LPush(ctx, key, k).Result()
}

func (r redis) SRem(key string, k string) (interface{}, error) {
	return r.Client.SRem(ctx, key, k).Result()
}
func (r redis) SIsMember(key string, k string) (bool, error) {
	return r.Client.SIsMember(ctx, key, k).Result()
}
func (r redis) SCard(key string) (int64, error) {
	return r.Client.SCard(ctx, key).Result()
}
func (r redis) SMembers(key string) (map[string]struct{}, error) {
	return r.Client.SMembersMap(ctx, key).Result()
}

func (r redis) SRandMember(key string) (interface{}, error) {
	return r.Client.SRandMember(ctx, key).Result()
}

func (r redis) LPop(key string) (interface{}, error) {
	return r.Client.LPop(ctx, key).Result()
}
func (r redis) LRem(key string, num int64, value interface{}) (interface{}, error) {
	return r.Client.LRem(ctx, key, num, value).Result()
}

// ZIncrBy 有序集合
func (r redis) ZIncrBy(key string, increment float64, member string) (float64, error) {
	return r.Client.ZIncrBy(ctx, key, increment, member).Result()
}

func (r redis) ZAdd(key string, increment float64, member string) (int64, error) {
	return r.Client.ZAdd(ctx, key, &redisV8.Z{
		Score:  increment,
		Member: member,
	}).Result()
}
func (r redis) ZRevRange(key string, start int64, end int64) ([]string, error) {
	return r.Client.ZRevRange(ctx, key, start, end).Result()
}

func (r redis) ZRange(key string, start int64, end int64) ([]string, error) {
	return r.Client.ZRange(ctx, key, start, end).Result()
}

func (r redis) ZScore(key string, Score string) (float64, error) {
	return r.Client.ZScore(ctx, key, Score).Result()
}

func (r redis) ZRevRank(key string, Score string) (int64, error) {
	return r.Client.ZRevRank(ctx, key, Score).Result()
}
func (r redis) SPop(key string) (string, error) {
	return r.Client.SPop(ctx, key).Result()
}

func (r redis) ZRem(key string, Str string) (int64, error) {
	return r.Client.ZRem(ctx, key, Str).Result()
}
func (r redis) SPops(key string, Num int64) ([]string, error) {
	return r.Client.SPopN(ctx, key, Num).Result()
}

func (r redis) HGetInt64(key string, field string) (int64, error) {
	return r.Client.HGet(ctx, key, field).Int64()
}
func (r redis) HLen(key string) (int64, error) {
	return r.Client.HLen(ctx, key).Result()
}

// GetBit 读取位图
func (r redis) GetBit(key string, UserId int64) (int64, error) {
	//redis.GetBit(context.Background(),keys,(userid-1)).Result()
	return r.Client.GetBit(context.Background(), key, UserId-1).Result()
}

// SetBit 写入位图1
func (r redis) SetBit(key string, UserId int64, Num int) (int64, error) {
	//res,err = redis.SetBit(context.Background(),keys,(userid-1),1).Result()
	return r.Client.SetBit(context.Background(), key, UserId-1, Num).Result()
}

// BitCount 获取点赞数量
func (r redis) BitCount(key string) (int64, error) {
	count := redisV8.BitCount{Start: 0, End: 1}
	return r.Client.BitCount(context.Background(), key, &count).Result()
}
