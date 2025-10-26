package main

import (
	"context"

	"github.com/kevin-chtw/tw_common/storage"
	"github.com/kevin-chtw/tw_proto/sproto"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
	"github.com/topfreegames/pitaya/v3/pkg/session"
	"google.golang.org/protobuf/types/known/anypb"
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

	msg := &sproto.NetStateReq{Uid: uid, Online: online}
	data, err := anypb.New(msg)
	if err != nil {
		logger.Log.Errorf("failed to marshal net state req: %v", err)
	}
	req := &sproto.MatchReq{
		Matchid: matching.MatchId,
		Req:     data,
	}
	rsp := &sproto.MatchAck{}
	app.RPCTo(context.Background(), matching.ServerId, matching.ServerType+".remote.message", rsp, req)
}
