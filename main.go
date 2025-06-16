package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kevin-chtw/tw_proxy_svr/service"
	"github.com/sirupsen/logrus"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/acceptor"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/config"
)

var app pitaya.Pitaya

func main() {
	serverType := "proxy"
	port := flag.Int("port", 3250, "port to listen on")
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	config := config.NewDefaultPitayaConfig()
	builder := pitaya.NewDefaultBuilder(true, serverType, pitaya.Cluster, map[string]string{}, *config)
	builder.AddAcceptor(acceptor.NewTCPAcceptor(fmt.Sprintf(":%d", *port)))
	app = builder.Build()

	defer app.Shutdown()

	logrus.Infof("Pitaya server of type %s started on port %d", serverType, *port)
	initServices()
	app.Start()
}

func initServices() {
	accountSvc := service.NewAccountSvc(app)
	app.Register(accountSvc, component.WithName("account"), component.WithNameFunc(strings.ToLower))
}
