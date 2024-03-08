package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/uzhenyu/framework/config"
	"time"
)

var RE *redis.Client

type T struct {
	Redis struct {
		Host string `json:"Host"`
		Port string `json:"Port"`
	} `json:"Redis"`
}

func withClient(dataID string, hand func(cli *redis.Client) error) error {
	newRedis := new(T)
	getConfig, err := config.GetConfig(dataID, "DEFAULT_GROUP")
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(getConfig), &newRedis)
	if err != nil {
		return err
	}
	RE = redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%:%v", newRedis.Redis.Host, newRedis.Redis.Port)})
	defer RE.Close()
	err = hand(RE)
	if err != nil {
		return err
	}
	return nil
}

func Exists(ctx context.Context, dataID string, key string) (bool, error) {
	var data int64
	var err error
	err = withClient(dataID, func(cli *redis.Client) error {
		data, err = cli.Exists(ctx, key).Result()
		return err
	})
	if err != nil {
		return false, err
	}
	if data > 0 {
		return true, nil
	}
	return false, nil
}

func GetRedis(ctx context.Context, dataID string, key string) (string, error) {
	var data string
	var err error
	err = withClient(dataID, func(cli *redis.Client) error {
		data, err = cli.Get(ctx, key).Result()
		return err
	})
	if err != nil {
		return "", err
	}
	return data, nil
}

func SetKey(ctx context.Context, dataID string, key string, val interface{}, duration time.Duration) error {
	return withClient(dataID, func(cli *redis.Client) error {
		return cli.Set(ctx, key, val, duration).Err()
	})
}
