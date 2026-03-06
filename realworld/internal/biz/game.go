package biz

import (
	"context"
	"time"

	gamev1 "realworld/api/game/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// MapBasic 领域层的地图基础信息模型，抽象自 data 层与 proto 层。
type MapBasic struct {
	ID           int64
	Name         string
	Description  string
	Width        int32
	Height       int32
	ThumbnailURL string
}

// MapDetail 领域层的地图详情模型，包含所有地块。
type MapDetail struct {
	ID     int64
	Name   string
	Width  int32
	Height int32
	Tiles  []MapTile
}

// MapTile 表示地图中的单个格子。
type MapTile struct {
	X    int32
	Y    int32
	Type gamev1.TileType
}

// GameSession 表示一局游戏对战的领域模型。
type GameSession struct {
	SessionID        string
	Mode             gamev1.GameMode
	MapID            int64
	LocalPlayerCount int32
	Status           gamev1.SessionStatus
	CreatedAt        time.Time
}

// GameResult 表示一局游戏结束后用于统计的结果数据。
type GameResult struct {
	SessionID      string
	UserID         uint64
	Mode           gamev1.GameMode
	MapID          int64
	IsWin          bool
	Player1Kills   int32
	Player2Kills   int32
	EnemyKills     int32
	DurationSecond int32
}

// MapRepo 定义地图相关的数据访问接口。
type MapRepo interface {
	ListMaps(ctx context.Context) ([]*MapBasic, error)
	GetMapDetail(ctx context.Context, id int64) (*MapDetail, error)
}

// GameSessionRepo 定义对局生命周期相关的数据访问接口。
type GameSessionRepo interface {
	CreateSession(ctx context.Context, s *GameSession) error
	UpdateSessionStatusAndTimes(ctx context.Context, sessionID string, status gamev1.SessionStatus, startedAt, finishedAt *time.Time) error
}

// GameStatsRepo 定义战绩记录相关的数据访问接口。
type GameStatsRepo interface {
	SaveMatchResult(ctx context.Context, r *GameResult) error
}

// GameConfigProvider 定义游戏静态配置（坦克、子弹等）的提供接口。
// difficulty: 1=简单 2=普通 3=困难，影响移速与子弹速度等。
type GameConfigProvider interface {
	GetTankConfig(ctx context.Context, difficulty int32) (*gamev1.TankConfigSet, error)
}

// GameUsecase 封装与游戏配置、地图与对局管理相关的业务用例。
type GameUsecase struct {
	maps     MapRepo
	sessions GameSessionRepo
	stats    GameStatsRepo
	config   GameConfigProvider
	log      *log.Helper
}

// NewGameUsecase 创建 GameUsecase 实例。
func NewGameUsecase(
	maps MapRepo,
	sessions GameSessionRepo,
	stats GameStatsRepo,
	config GameConfigProvider,
	logger log.Logger,
) *GameUsecase {
	return &GameUsecase{
		maps:     maps,
		sessions: sessions,
		stats:    stats,
		config:   config,
		log:      log.NewHelper(logger),
	}
}

// ListMaps 返回所有可用地图的基础信息。
func (uc *GameUsecase) ListMaps(ctx context.Context) ([]*MapBasic, error) {
	uc.log.Infof("list all maps")
	return uc.maps.ListMaps(ctx)
}

// GetMap 返回指定 ID 的地图详情。
func (uc *GameUsecase) GetMap(ctx context.Context, id int64) (*MapDetail, error) {
	uc.log.Infof("get map detail, id=%d", id)
	return uc.maps.GetMapDetail(ctx, id)
}

// GetTankConfig 根据难度返回坦克与子弹的基础属性配置。
func (uc *GameUsecase) GetTankConfig(ctx context.Context, difficulty int32) (*gamev1.TankConfigSet, error) {
	uc.log.Infof("get tank config set, difficulty=%d", difficulty)
	return uc.config.GetTankConfig(ctx, difficulty)
}

// CreateSession 创建一局新的游戏对局，并返回业务上的 sessionID。
func (uc *GameUsecase) CreateSession(ctx context.Context, mode gamev1.GameMode, mapID int64, localPlayerCount int32) (string, error) {
	uc.log.Infof("create game session, mode=%v, mapID=%d, players=%d", mode, mapID, localPlayerCount)

	sessionID := generateSessionID()
	now := time.Now()
	session := &GameSession{
		SessionID:        sessionID,
		Mode:             mode,
		MapID:            mapID,
		LocalPlayerCount: localPlayerCount,
		Status:           gamev1.SessionStatus_SESSION_STATUS_CREATED,
		CreatedAt:        now,
	}

	if err := uc.sessions.CreateSession(ctx, session); err != nil {
		uc.log.Errorf("create game session failed: %v", err)
		return "", err
	}

	return sessionID, nil
}

// FinishSession 在对局结束时记录结果并更新对局状态。
func (uc *GameUsecase) FinishSession(ctx context.Context, res *GameResult) error {
	uc.log.Infof("finish game session, id=%s, win=%v", res.SessionID, res.IsWin)

	now := time.Now()
	if err := uc.sessions.UpdateSessionStatusAndTimes(
		ctx,
		res.SessionID,
		gamev1.SessionStatus_SESSION_STATUS_FINISHED,
		nil,
		&now,
	); err != nil {
		uc.log.Errorf("update session status failed: %v", err)
		return err
	}

	if err := uc.stats.SaveMatchResult(ctx, res); err != nil {
		uc.log.Errorf("save match result failed: %v", err)
		return err
	}

	return nil
}

// generateSessionID 生成一个简单的 sessionID。
// 初版使用时间戳+随机数，后续可替换为更规范的 UUID。
func generateSessionID() string {
	return time.Now().Format("20060102150405.000000000")
}
