package logic

import (
	"context"
	"database/sql"
	"time"

	"gim/apps/social/models"
	"gim/apps/social/rpc/internal/svc"
	"gim/apps/social/rpc/social"
	"gim/pkg/constants"
	"gim/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line

	friends,err:=l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx,in.UserId,in.ReqUid)
	if err != nil && err != models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v ", err, in)
	}
	if friends != nil {
		return &social.FriendPutInResp{}, err
	}
	
	friendReqs,err:=l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx,in.ReqUid,in.UserId)
	if err != nil && err != models.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err %v req %v ", err, in)
	}
	if friendReqs != nil {
		return &social.FriendPutInResp{}, err
	}

	_,err=l.svcCtx.FriendRequestsModel.Insert(l.ctx, &models.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}
	return &social.FriendPutInResp{}, nil
}
