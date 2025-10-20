package main

import (
	"strings"

	"github.com/kevin-chtw/tw_common/storage"
	"github.com/kevin-chtw/tw_common/utils"
	"github.com/kevin-chtw/tw_proxy_svr/service"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/serialize"
)

var app pitaya.Pitaya

func main() {
	serverType := "proxy"
	pitaya.SetLogger(utils.Logger(logrus.DebugLevel))

	config := config.NewDefaultPitayaConfig()
	config.SerializerType = uint16(serialize.PROTOBUF)
	config.Handler.Messages.Compression = false
	// config.Cluster.RPC.Client.Nats.Connect = "nats://192.168.182.128:4222"
	// config.Cluster.RPC.Server.Nats.Connect = "nats://192.168.182.128:4222"
	// config.Cluster.SD.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	// config.Groups.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(acceptor.NewTCPAcceptor(":3250"))
	builder.AddAcceptor(acceptor.NewWSAcceptor(":3251"))
	builder.Router.AddRoute("game", GameRouter)
	builder.SessionPool.OnSessionClose(OnSessionClose)
	builder.SessionPool.OnAfterSessionBind(OnAfterSessionBind)
	app = builder.Build()

	defer app.Shutdown()
	bs := storage.NewETCDMatching(builder.Server, builder.Config.Modules.BindingStorage.Etcd)
	app.RegisterModule(bs, "matchingstorage")

	initServices()
	app.Start()
}

func initServices() {
	player := service.NewPlayer(app)
	app.Register(player, component.WithName("player"), component.WithNameFunc(strings.ToLower))
}
