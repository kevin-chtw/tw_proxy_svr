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
	go sendNetState(uid, false)
}

func OnAfterSessionBind(ctx context.Context, s session.Session) error {
	uid := s.UID()
	logger.Log.Infof("session binded: %s", uid)
	go sendNetState(uid, true)
	return nil
}

func sendNetState(uid string, online bool) {
	module, err := app.GetModule("matchingstorage")
	if err != nil {
		logger.Log.Errorf("get module error: %v", err)
		return
	}
	ms := module.(*storage.ETCDMatching)
	matching, err := ms.Get(uid)
	if err != nil || matching == nil {
		logger.Log.Warnf("get matching error: %v", err)
		return
	}

	req := &sproto.NetStateReq{Uid: uid, Online: online}
	rsp := &sproto.NetStateAck{}
	app.RPCTo(context.Background(), matching.ServerId, matching.ServerType+".player.net", rsp, req)
}
