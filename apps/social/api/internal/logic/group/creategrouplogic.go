package group

import (
	"context"

	"gim/apps/im/rpc/imclient"
	"gim/apps/social/api/internal/svc"
	"gim/apps/social/api/internal/types"
	"gim/apps/social/rpc/socialclient"
	"gim/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创群
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// todo: add your logic here and delete this line

	uid:=ctxdata.GetUId(l.ctx)

	res,err:=l.svcCtx.Social.GroupCreate(l.ctx,&socialclient.GroupCreateReq{
		Name: req.Name,
		Icon: req.Icon,
		CreatorUid: uid,
	})

	if err!=nil{
		return nil, err
	}

	if res.Id==""{
		return nil,err
	}

	l.svcCtx.Im.CreateGroupConversation(l.ctx,&imclient.CreateGroupConversationReq{
		GroupId: res.Id,
		CreateId: uid,
	})
	return nil,err
}
