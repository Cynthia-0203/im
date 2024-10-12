package svc

import (
	"gim/apps/user/api/internal/config"
	"gim/apps/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "user.User"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN"]
		}
	}]
}`


type ServiceContext struct {
	Config config.Config

	userclient.User
	Redis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc,zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)))),
		Redis: redis.MustNewRedis(c.Redisx),
	}
}
