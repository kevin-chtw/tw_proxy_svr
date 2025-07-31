package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/config"
)

var app pitaya.Pitaya

func main() {
	serverType := "proxy"
	port := flag.Int("port", 3250, "port to listen on")
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	config := config.NewDefaultPitayaConfig()
	config.Cluster.RPC.Client.Nats.Connect = "nats://192.168.182.128:4222"
	config.Cluster.RPC.Server.Nats.Connect = "nats://192.168.182.128:4222"
	config.Cluster.SD.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	config.Groups.Etcd.Endpoints = []string{"http://192.168.182.128:2379"}
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port)))
	app = builder.Build()

	defer app.Shutdown()

	app.Start()
}
