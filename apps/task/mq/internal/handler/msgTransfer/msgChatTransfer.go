package msgtransfer

import (
	"context"
	"encoding/json"
	model "gim/apps/im/models"

	"gim/apps/im/ws/ws"

	"gim/apps/task/mq/internal/svc"
	"gim/apps/task/mq/mq"
	"gim/pkg/bitmap"


	// "github.com/zeromicro/go-queue/kq"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context,key, value string) error {
	var (
		data  mq.MsgChatTransfer
		// ctx   = context.Background()
		msgId = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := model.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
