package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/antgubarev/pingbot/internal/model"
	"github.com/go-redis/redis/v9"
)

const (
	keyTemplate = "response_time:%s" // user_id:target_id
)

type RedisResponseStorage struct {
	rdb *redis.Client
}

type RedisOpt struct {
	Host     string
	Port     string
	Password *string
}

func ParseRedisOpt() (*RedisOpt, error) {
	opt := RedisOpt{}
	opt.Host = os.Getenv("REDIS_HOST")
	opt.Port = os.Getenv("REDIS_PORT")
	if pass, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		opt.Password = &pass
	}

	if opt.Host == "" || opt.Port == "" {
		return nil, fmt.Errorf("REDIS_HOST and REDIS_PORT must be both set")
	}

	return &opt, nil
}

func NewRedisResponseStorage(opt *RedisOpt) *RedisResponseStorage {
	redisOpt := redis.Options{
		Addr: fmt.Sprintf("%s:%s", opt.Host, opt.Port),
	}
	if opt.Password != nil {
		redisOpt.Password = *opt.Password
	}
	rdb := redis.NewClient(&redisOpt)
	return &RedisResponseStorage{
		rdb: rdb,
	}
}

func (rrs *RedisResponseStorage) Save(ctx context.Context, response *model.ResponseTime, target *model.Target) error {
	_, err := rrs.rdb.ZAdd(ctx, fmt.Sprintf(keyTemplate, target.Id), redis.Z{
		Score:  float64(response.RTime.Unix()),
		Member: strconv.FormatInt(response.DurationMs, 10),
	}).Result()

	if err != nil {
		return fmt.Errorf("ZAdd response data: %v", err)
	}

	return nil
}

func (rrs *RedisResponseStorage) GetLastResonseTime(ctx context.Context, targetId string) (*time.Time, error) {
	items, err := rrs.rdb.ZRangeByScoreWithScores(ctx, fmt.Sprintf(keyTemplate, targetId), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()

	if err != nil {

		return nil, fmt.Errorf("ZRangeByScore: GetLastResonseTime: %v", err)
	}
	if len(items) == 0 {

		return nil, nil
	}
	rtime := time.Unix(int64(items[0].Score), 0)

	return &rtime, nil
}

func (rrs *RedisResponseStorage) GetResponseTimesByPeriod(ctx context.Context, from, to time.Time, targetId string) ([]model.ResponseTime, error) {
	items, err := rrs.rdb.ZRangeByScoreWithScores(ctx, fmt.Sprintf(keyTemplate, targetId), &redis.ZRangeBy{
		Min:    strconv.FormatInt(from.Unix(), 10),
		Max:    strconv.FormatInt(to.Unix(), 10),
		Offset: 0,
		Count:  10000,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("ZRange: GetResponseTimesByPeriod: %v", err)
	}

	result := make([]model.ResponseTime, len(items))
	for index, item := range items {
		rtime := time.Unix(int64(item.Score), 0)
		durationString, ok := item.Member.(string)
		if !ok {
			//TODO: logging
			result[index] = model.ResponseTime{
				DurationMs: -1,
				RTime:      rtime,
			}
			continue
		}
		duration, err := strconv.ParseInt(durationString, 10, 64)
		if err != nil {
			//TODO: logging
			result[index] = model.ResponseTime{
				DurationMs: -1,
				RTime:      rtime,
			}
			continue
		}
		result[index] = model.ResponseTime{
			DurationMs: duration,
			RTime:      rtime,
		}
	}

	return result, nil
}
