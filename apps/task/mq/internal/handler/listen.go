package handler

import (
	
	"gim/apps/task/mq/internal/svc"
	"gim/apps/task/mq/internal/handler/msgTransfer"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Listen struct{
	svc  *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext)*Listen{
	return &Listen{
		svc: svc,
	}
}

func (l *Listen) Services() []service.Service {
	return []service.Service{
		// todo: 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer,msgtransfer.NewMsgChatTransfer(l.svc)),
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer,msgtransfer.NewMsgReadTransfer(l.svc)),
	}
}
