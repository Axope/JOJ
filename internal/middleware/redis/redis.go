package redis

import (
	"context"

	"github.com/Axope/JOJ/configs"
	"github.com/redis/go-redis/v9"
)

type Z = redis.Z
type Tx = redis.Tx
type Pipeliner = redis.Pipeliner
type StringCmd = redis.StringCmd

const (
	TxFailedErr = redis.TxFailedErr
)

var (
	ctx = context.Background()

	cfg configs.RedisConfig
	rdb *redis.Client
)

func InitRedis() error {
	cfg = configs.GetRedisConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

// func ZAdd(key string, score float64, uid uint) error {
// 	return rdb.ZAdd(ctx, key, redis.Z{
// 		Score:  score,
// 		Member: uid,
// 	}).Err()
// }
// func ZRem(key string, uid uint) error {
// 	return rdb.ZRem(ctx, key, uid).Err()
// }
func ZRangeWithScores(ctx context.Context, key string, startIdx, stopIdx int64) ([]redis.Z, error) {
	return rdb.ZRangeWithScores(ctx, key, startIdx, stopIdx).Result()
}
// func ZCard(key string) (int64, error) {
// 	return rdb.ZCard(ctx, key).Result()
// }

//	func Set(key string, value interface{}) error {
//		return rdb.Set(ctx, key, value, 0).Err()
//	}
func Del(ctx context.Context, key string) error {
	return rdb.Del(ctx, key).Err()
}

// func Get(key string) (string, error) {
// 	return rdb.Get(ctx, key).Result()
// }

func Watch(ctx context.Context, txf func(*redis.Tx) error, keys ...string) error {
	return rdb.Watch(ctx, txf, keys...)
}

func Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return rdb.Pipelined(ctx, fn)
}

func Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return rdb.Scan(ctx, cursor, match, count)
}
