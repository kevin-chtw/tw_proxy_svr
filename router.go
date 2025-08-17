package main

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/kevin-chtw/tw_proto/cproto"
	"github.com/topfreegames/pitaya/v3/pkg/cluster"
	"github.com/topfreegames/pitaya/v3/pkg/route"
)

// 自定义 Router
func GameRouter(ctx context.Context, route *route.Route, payload []byte, servers map[string]*cluster.Server) (*cluster.Server, error) {
	req := &cproto.GameReq{}
	if err := json.Unmarshal(payload, &req); err == nil && req.Serverid != "" {
		if svr, ok := servers[req.Serverid]; ok {
			return svr, nil
		}
	}

	srvList := make([]*cluster.Server, 0)
	for _, v := range servers {
		srvList = append(srvList, v)
	}
	server := srvList[rand.Intn(len(srvList))]
	return server, nil
}
