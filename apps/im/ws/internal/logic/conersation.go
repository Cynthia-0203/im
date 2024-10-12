package logic

import (
	"context"
	"time"

	model "gim/apps/im/models"
	"gim/apps/im/ws/internal/svc"
	"gim/apps/im/ws/websocket"
	"gim/apps/im/ws/ws"
	"gim/pkg/wuid"
)

// type ChatLogSlg interface{
// 	SingleChatLog(data *types.Chat,userId string)error
// }

type UserLogic struct{
	ctx context.Context
	srv *websocket.Server
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context,srv *websocket.Server,svcCtx *svc.ServiceContext)*UserLogic{
	return &UserLogic{
		ctx: ctx,
		srv: srv,
		svcCtx: svcCtx,
	}
}

func(l *UserLogic)Chat(data *ws.Chat,userId string)error{
	if data.ConversationId==""{
		data.ConversationId=wuid.CombineId(userId,data.RecvId)
	}

	// time.Sleep(time.Minute)

	chatLog:=model.ChatLog{
		ConversationId: data.ConversationId,
		SendId: userId,
		RecvId: data.RecvId,
		ChatType: data.ChatType,
		MsgFrom: 0,
		MsgType: data.MType,
		MsgContent: data.Content,
		SendTime: time.Now().UnixNano(),
	}

	err:=l.svcCtx.ChatLogModel.Insert(l.ctx,&chatLog)

	return err
}