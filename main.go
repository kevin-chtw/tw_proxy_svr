package main

import (
	"flag"
	"fmt"

	"github.com/kevin-chtw/tw_common/storage"
	"github.com/kevin-chtw/tw_common/utils"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/config"
	"github.com/topfreegames/pitaya/v3/pkg/serialize"
)

var app pitaya.Pitaya

func main() {
	serverType := "proxy"
	pitaya.SetLogger(utils.Logger(logrus.DebugLevel))
	port := flag.Int("port", 3250, "port to listen on")
	flag.Parse()

	config := config.NewDefaultPitayaConfig()
	config.SerializerType = uint16(serialize.PROTOBUF)
	config.Handler.Messages.Compression = false
	// config.Cluster.RPC.Client.Nats.Connect = "nats://192.168.182.128:4222"
	// config.Cluster.RPC.Server.Nats.Connect = "nats://192.168.182.128:4222"
	// config.Cluster.SD.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	// config.Groups.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port)))
	builder.Router.AddRoute("game", GameRouter)
	builder.SessionPool.OnSessionClose(OnSessionClose)
	builder.SessionPool.OnAfterSessionBind(OnAfterSessionBind)
	app = builder.Build()
	defer app.Shutdown()
	bs := storage.NewETCDMatching(builder.Server, builder.Config.Modules.BindingStorage.Etcd)
	app.RegisterModule(bs, "matchingstorage")

	app.Start()
}
