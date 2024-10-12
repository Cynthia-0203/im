package svc

import (
	"gim/apps/user/models"
	"gim/apps/user/rpc/internal/config"
	"gim/pkg/constants"
	"gim/pkg/ctxdata"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		Redis:redis.MustNewRedis(c.Redisx),
		UsersModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func(svcCtx *ServiceContext)SetRootToken()error{
	systemToken,err:=ctxdata.GetJwtToken(svcCtx.Config.Jwt.AccessSecret,time.Now().Unix(),9999999,constants.SYSTEM_ROOT_UID)
	if err!=nil{
		return err
	}

	return svcCtx.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN,systemToken)
}