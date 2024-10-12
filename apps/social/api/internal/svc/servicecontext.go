package svc

import (
	"gim/apps/im/rpc/imclient"
	"gim/apps/social/api/internal/config"
	"time"

	"gim/apps/social/rpc/socialclient"
	"gim/apps/user/rpc/userclient"

	// "gim/pkg/interceptor"
	"gim/pkg/interceptor"
	"gim/pkg/interceptor/rpcclient"
	"gim/pkg/middleware"

	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	// "google.golang.org/grpc"
)

var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

type ServiceContext struct {
	Config                config.Config
	IdempotenceMiddleware rest.Middleware
	LimitMiddleware       rest.Middleware
	*redis.Redis
	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Im: imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		Redis:                 redis.MustNewRedis(c.Redisx),
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc,
			zrpc.WithUnaryClientInterceptor(rpcclient.NewSheddingClient("user-rpc",
				load.WithBuckets(10),
				load.WithCpuThreshold(1),
				load.WithWindow(time.Millisecond*100000),
			)),
		)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient),
		)),

		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		LimitMiddleware:       middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(1, 100),
	}
}
