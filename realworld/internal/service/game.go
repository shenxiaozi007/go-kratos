package service

import (
	"context"

	gamev1 "realworld/api/game/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// GameService 实现 game.v1.GameServiceServer 与 game.v1.GameServiceHTTPServer 接口，
// 负责将 HTTP/gRPC 请求转换为 GameUsecase 的调用，并组装返回结果。
type GameService struct {
	gamev1.UnimplementedGameServiceServer

	uc  *biz.GameUsecase
	log *log.Helper
}

// NewGameService 创建 GameService 实例。
func NewGameService(uc *biz.GameUsecase, logger log.Logger) *GameService {
	return &GameService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ListMaps 获取所有可选地图列表。
func (s *GameService) ListMaps(ctx context.Context, _ *gamev1.ListMapsRequest) (*gamev1.ListMapsReply, error) {
	maps, err := s.uc.ListMaps(ctx)
	if err != nil {
		s.log.Errorf("ListMaps failed: %v", err)
		return nil, err
	}

	reply := &gamev1.ListMapsReply{
		Maps: make([]*gamev1.MapBasicInfo, 0, len(maps)),
	}
	for _, m := range maps {
		reply.Maps = append(reply.Maps, &gamev1.MapBasicInfo{
			Id:           m.ID,
			Name:         m.Name,
			Description:  m.Description,
			Width:        m.Width,
			Height:       m.Height,
			ThumbnailUrl: m.ThumbnailURL,
		})
	}
	return reply, nil
}

// GetMap 获取单张地图详细配置。
func (s *GameService) GetMap(ctx context.Context, req *gamev1.GetMapRequest) (*gamev1.GetMapReply, error) {
	detail, err := s.uc.GetMap(ctx, req.GetId())
	if err != nil {
		s.log.Errorf("GetMap failed: %v", err)
		return nil, err
	}

	reply := &gamev1.GetMapReply{
		Map: &gamev1.MapDetail{
			Id:     detail.ID,
			Name:   detail.Name,
			Width:  detail.Width,
			Height: detail.Height,
			Tiles:  make([]*gamev1.MapTile, 0, len(detail.Tiles)),
		},
	}
	for _, t := range detail.Tiles {
		reply.Map.Tiles = append(reply.Map.Tiles, &gamev1.MapTile{
			X:    t.X,
			Y:    t.Y,
			Type: t.Type,
		})
	}
	return reply, nil
}

// GetTankConfig 根据难度获取坦克与子弹的基础属性配置。
func (s *GameService) GetTankConfig(ctx context.Context, req *gamev1.GetTankConfigRequest) (*gamev1.GetTankConfigReply, error) {
	difficulty := req.GetDifficulty()
	if difficulty == 0 {
		difficulty = 2 // 默认普通
	}
	cfg, err := s.uc.GetTankConfig(ctx, difficulty)
	if err != nil {
		s.log.Errorf("GetTankConfig failed: %v", err)
		return nil, err
	}
	return &gamev1.GetTankConfigReply{Config: cfg}, nil
}

// CreateSession 创建对局。
func (s *GameService) CreateSession(ctx context.Context, req *gamev1.CreateSessionRequest) (*gamev1.CreateSessionReply, error) {
	sessionID, err := s.uc.CreateSession(ctx, req.GetMode(), req.GetMapId(), req.GetLocalPlayerCount())
	if err != nil {
		s.log.Errorf("CreateSession failed: %v", err)
		return nil, err
	}
	return &gamev1.CreateSessionReply{SessionId: sessionID}, nil
}

// FinishSession 上报对局结束结果。
func (s *GameService) FinishSession(ctx context.Context, req *gamev1.FinishSessionRequest) (*gamev1.FinishSessionReply, error) {
	res := &biz.GameResult{
		SessionID:      req.GetSessionId(),
		UserID:         0, // 初版使用游客模式，后续可从登录态中获取真实 userID。
		Mode:           req.GetMode(),
		MapID:          req.GetMapId(),
		IsWin:          req.GetIsWin(),
		Player1Kills:   req.GetPlayer1Kills(),
		Player2Kills:   req.GetPlayer2Kills(),
		EnemyKills:     req.GetEnemyKills(),
		DurationSecond: req.GetDurationSeconds(),
	}

	if err := s.uc.FinishSession(ctx, res); err != nil {
		s.log.Errorf("FinishSession failed: %v", err)
		return nil, err
	}

	return &gamev1.FinishSessionReply{
		SessionId: req.GetSessionId(),
		Status:    gamev1.SessionStatus_SESSION_STATUS_FINISHED,
	}, nil
}

