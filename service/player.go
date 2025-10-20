package service

import (
	"context"
	"encoding/json"

	"github.com/kevin-chtw/tw_common/utils"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/logger"
)

type ClientInfo struct {
	NetType string `json:"net_type"`
}

type Player struct {
	component.Base
	app pitaya.Pitaya
}

func NewPlayer(app pitaya.Pitaya) *Player {
	return &Player{
		app: app,
	}
}

func (p *Player) Message(ctx context.Context, data []byte) {
	req := &ClientInfo{}
	if err := json.Unmarshal(data, req); err != nil {
		logger.Log.Error(err.Error())
		return
	}

	s := p.app.GetSessionFromCtx(ctx)
	if err := s.Set(utils.NetType, req.NetType); err != nil {
		logger.Log.Error(err.Error())
	}
}
