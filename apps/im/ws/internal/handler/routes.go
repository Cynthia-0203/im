package handler

import (
	"gim/apps/im/ws/internal/handler/conversation"
	"gim/apps/im/ws/internal/handler/push"
	"gim/apps/im/ws/internal/handler/user"
	"gim/apps/im/ws/internal/svc"
	"gim/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server,svc *svc.ServiceContext){
	srv.AddRoutes([]websocket.Route{
		{
			Method: "user.online",
			Handler: user.OnLine(svc),
		},
		{
			Method: "conversation.chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method: "conversation.markChat",
			Handler: conversation.MarkRead(svc),
		},
		{
			Method: "push",
			Handler: push.Push(svc),
		},
	})
}

