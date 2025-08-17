package main

import (
	"context"

	"github.com/kevin-chtw/tw_proto/sproto"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
	"github.com/topfreegames/pitaya/v3/pkg/modules"
	"github.com/topfreegames/pitaya/v3/pkg/session"
)

func OnSessionClose(s session.Session) {
	uid := s.UID()
	logger.Log.Infof("session closed: %s", uid)
	module, err := app.GetModule("matchingstorage")
	if err != nil {
		return
	}
	ms := module.(*modules.ETCDBindingStorage)
	serverId, err := ms.GetUserFrontendID(uid, "match")
	if err != nil || serverId == "" {
		return
	}

	req := &sproto.Proxy2MatchReq{
		Req: &sproto.Proxy2MatchReq_OfflineReq{
			OfflineReq: &sproto.OfflineReq{
				Uid: uid,
			},
		},
	}
	rsp := &sproto.Proxy2MatchAck{}
	app.RPCTo(context.Background(), serverId, "match.player.offline", rsp, req)
}

func OnAfterSessionBind(ctx context.Context, s session.Session) error {
	uid := s.UID()
	logger.Log.Infof("session binded: %s", uid)
	module, err := app.GetModule("matchingstorage")
	if err != nil {
		return nil
	}
	ms := module.(*modules.ETCDBindingStorage)
	serverId, err := ms.GetUserFrontendID(uid, "match")
	if err != nil || serverId == "" {
		return nil
	}

	req := &sproto.Proxy2MatchReq{
		Req: &sproto.Proxy2MatchReq_OnlineReq{
			OnlineReq: &sproto.OnlineReq{
				Uid: uid,
			},
		},
	}
	rsp := &sproto.Proxy2MatchAck{}
	if err := app.RPCTo(ctx, serverId, "match.player.online", rsp, req); err != nil {
		logger.Log.Error(err)
	}
	return nil
}
