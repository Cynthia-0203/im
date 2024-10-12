package svc

import (
	model "gim/apps/im/models"
	"gim/apps/im/ws/websocket"
	"gim/apps/social/rpc/socialclient"
	"gim/apps/task/mq/internal/config"
	"gim/pkg/constants"
	"net/http"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	config.Config

	WsClient websocket.Client
	*redis.Redis

	socialclient.Social
	model.ChatLogModel
	model.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      model.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: model.MustConversationModel(c.Mongo.Url, c.Mongo.Db),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(
		c.Ws.Host,
		websocket.WithClientHeader(header),
		websocket.WithClientDiscover(websocket.NewRedisDiscover(header, constants.REDIS_DISCOVER_SRV, c.Redisx)),
	)
	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
