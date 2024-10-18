package rds

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"github.com/tigapilarmandiri/perkakas/configs"
)

var (
	rediser   Rediser
	onceRedis sync.Once
)

func GetClient() Rediser {
	onceRedis.Do(func() {
		if configs.Config.Redis.IsCluster {
			// if it redis cluster
			clusterServiceAddress := configs.Config.Redis.Host + ":" + configs.Config.Redis.Port
			clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:          []string{clusterServiceAddress},
				Username:       configs.Config.Redis.User,
				Password:       configs.Config.Redis.Pass,
				PoolSize:       configs.Config.Redis.PoolSize,
				MaxActiveConns: configs.Config.Redis.MaxActiveConns,
			})

			if configs.Config.Redis.Enabled {
				err := clusterClient.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
					return shard.Ping(ctx).Err()
				})
				if err != nil {
					util.Log.Fatal().Msg(err.Error())
				}
			}

			rediser = clusterClient

			return
		}

		// if it non cluster
		client := redis.NewClient(&redis.Options{
			Addr:           configs.Config.Redis.Host + ":" + configs.Config.Redis.Port,
			Username:       configs.Config.Redis.User,
			Password:       configs.Config.Redis.Pass,
			DB:             configs.Config.Redis.DB,
			PoolSize:       configs.Config.Redis.PoolSize,
			MaxActiveConns: configs.Config.Redis.MaxActiveConns,
		})

		if configs.Config.Redis.Enabled {
			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				util.Log.Fatal().Msg(err.Error())
			}
		}

		rediser = client
	})
	return rediser
}
