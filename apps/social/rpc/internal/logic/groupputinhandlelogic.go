package logic

import (
	"context"
	"database/sql"

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
	ErrGroupReqBeforePass   = xerr.NewMsg("Request passed")
	ErrGroupReqBeforeRefuse = xerr.NewMsg("Request rejected")
)
type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// todo: add your logic here and delete this line

	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend req err %v req %v", err, in.GroupReqId)
	}

	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}

	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}

	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx,  groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend req err %v req %v", err, groupReq)
		}

		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}

		groupMember := &models.GroupMembers{
			GroupId:     groupReq.GroupId,
			UserId:      groupReq.ReqId,
			RoleLevel:   int64(constants.AtLargeGroupRoleLevel),
			OperatorUid: in.HandleUid,
		}
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx,  groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
		}

		return nil
	})

	if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, err
}
