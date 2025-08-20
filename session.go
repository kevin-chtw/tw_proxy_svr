package main

import (
	"context"

	"github.com/kevin-chtw/tw_common/storage"
	"github.com/kevin-chtw/tw_proto/sproto"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
	"github.com/topfreegames/pitaya/v3/pkg/session"
)

func OnSessionClose(s session.Session) {
	uid := s.UID()
	logger.Log.Infof("session closed: %s", uid)
	module, err := app.GetModule("matchingstorage")
	if err != nil {
		return
	}
	ms := module.(*storage.ETCDMatching)
	matching, err := ms.Get(uid)
	if err != nil || matching == nil {
		return
	}

	req := &sproto.NetStateReq{Uid: uid, Online: false}
	rsp := &sproto.NetStateAck{}
	app.RPCTo(context.Background(), matching.ServerId, matching.ServerType+".player.net", rsp, req)
}

func OnAfterSessionBind(ctx context.Context, s session.Session) error {
	uid := s.UID()
	logger.Log.Infof("session binded: %s", uid)
	module, err := app.GetModule("matchingstorage")
	if err != nil {
		return nil
	}
	ms := module.(*storage.ETCDMatching)
	matching, err := ms.Get(uid)
	if err != nil || matching == nil {
		return nil
	}

	req := &sproto.NetStateReq{Uid: uid, Online: true}
	rsp := &sproto.NetStateAck{}
	if err := app.RPCTo(ctx, matching.ServerId, matching.ServerType+".player.session", rsp, req); err != nil {
		logger.Log.Error(err)
	}
	return nil
}
