package svc

import (
	model "gim/apps/im/models"
	"gim/apps/im/ws/internal/config"
	"gim/apps/task/mq/mqclient"
)


type ServiceContext struct{
	Config config.Config
	model.ChatLogModel

	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
}

func NewServiceContext(c config.Config)*ServiceContext{
	return &ServiceContext{
		Config: c,

		ChatLogModel: model.MustChatLogModel(c.Mongo.Url,c.Mongo.Db),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
	}
}

