package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kevin-chtw/tw_common/utils"
	"github.com/kevin-chtw/tw_proto/sproto"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
	"github.com/topfreegames/pitaya/v3/pkg/modules"
	"github.com/topfreegames/pitaya/v3/pkg/session"
)

var app pitaya.Pitaya

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

func main() {
	serverType := "proxy"
	pitaya.SetLogger(utils.Logger(logrus.DebugLevel))
	port := flag.Int("port", 3250, "port to listen on")
	flag.Parse()

	config := config.NewDefaultPitayaConfig()
	// config.Cluster.RPC.Client.Nats.Connect = "nats://192.168.182.128:4222"
	// config.Cluster.RPC.Server.Nats.Connect = "nats://192.168.182.128:4222"
	// config.Cluster.SD.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	// config.Groups.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port)))
	builder.SessionPool.OnSessionClose(OnSessionClose)
	builder.SessionPool.OnAfterSessionBind(OnAfterSessionBind)
	app = builder.Build()
	defer app.Shutdown()
	bs := modules.NewETCDBindingStorage(builder.Server, builder.SessionPool, builder.Config.Modules.BindingStorage.Etcd)
	app.RegisterModule(bs, "matchingstorage")

	app.Start()
}
