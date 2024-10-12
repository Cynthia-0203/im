package handler

import (
	"net/http"

	"gim/apps/im/api/internal/logic"
	"gim/apps/im/api/internal/svc"
	"gim/apps/im/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取会话已读记录
func GetChatLogReadRecordsReqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetChatLogReadRecordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetChatLogReadRecordsReqLogic(r.Context(), svcCtx)
		resp, err := l.GetChatLogReadRecordsReq(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
