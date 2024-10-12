package logic

import (
	"context"

	"gim/apps/social/models"
	"gim/apps/social/rpc/internal/svc"
	"gim/apps/social/rpc/social"
	"gim/pkg/constants"
	"gim/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)


var (
	ErrFriendReqBeforePass   = xerr.NewMsg("Friend application has been approved")
	ErrFriendReqBeforeRefuse = xerr.NewMsg("Friend request has been rejected")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line
	friendReq,err:=l.svcCtx.FriendRequestsModel.FindOne(l.ctx,uint64(in.FriendReqId))
	if err!=nil{
		return nil,errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err %v req %v ", err,in.FriendReqId)
	}

	switch constants.HandlerResult(friendReq.HandleResult.Int64){
	case constants.PassHandlerResult:
		return nil,errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil,errors.WithStack(ErrFriendReqBeforeRefuse)
	}

	friendReq.HandleResult.Int64=int64(in.HandleResult)

	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, friendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v, req %v", err, friendReq)
		}

		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		friends := []*models.Friends{
			{
				UserId:    friendReq.UserId,
				FriendUid: friendReq.ReqUid,
			}, {
				UserId:    friendReq.ReqUid,
				FriendUid: friendReq.UserId,
			},
		}

		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "friends inserts err %v, req %v", err, friends)
		}
		return nil
	})

	return &social.FriendPutInHandleResp{}, nil
}
