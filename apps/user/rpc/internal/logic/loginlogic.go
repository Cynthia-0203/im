package logic

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"gim/apps/user/models"
	"gim/apps/user/rpc/internal/svc"
	"gim/apps/user/rpc/user"
	"gim/pkg/ctxdata"
	"gim/pkg/encrypt"
	"gim/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var(
	ErrPhoneNotRegister=xerr.New(xerr.SERVER_COMMON_ERROR,"this phone_number is not registered")
	ErrUserPwdError =xerr.New(xerr.SERVER_COMMON_ERROR,"incorrect password")
)
type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line
	userEntity,err:=l.svcCtx.UsersModel.FindByPhone(l.ctx,in.Phone)

	if err!=nil{
		if err==models.ErrNotFound{
			return nil,errors.WithStack(ErrPhoneNotRegister)
		}
		return nil,errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}
	if !encrypt.ValidatePasswordHash(in.Password,userEntity.Password.String){
		return nil,errors.WithStack(ErrUserPwdError)
	}

	now:=time.Now().Unix()
	token,err:=ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret,now,l.svcCtx.Config.Jwt.AccessExpire,userEntity.Id)

	if err!=nil{
		return nil,errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	
	return &user.LoginResp{
		Id: userEntity.Id,
		Token: token,
		Expire: now+l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
