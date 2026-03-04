package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// MatchRecordView 单条战绩展示用结构（与 stats.v1.MatchRecord 对齐）。
type MatchRecordView struct {
	SessionID      string
	UserID         uint64
	Mode           int32
	MapID          int64
	IsWin          bool
	Player1Kills   int32
	Player2Kills   int32
	EnemyKills     int32
	DurationSecond int32
	CreatedAt      int64 // unix 秒
}

// LeaderboardEntryView 排行榜一条记录。
type LeaderboardEntryView struct {
	UserID       uint64
	Username     string
	TotalMatches int64
	TotalWins    int64
	WinRate      float32
}

// StatsRepo 战绩与排行榜数据访问接口。
type StatsRepo interface {
	ListMatchRecords(ctx context.Context, page, pageSize int32, userID *uint64) ([]*MatchRecordView, int64, error)
	GetLeaderboard(ctx context.Context, limit int32) ([]*LeaderboardEntryView, error)
}

// StatsUsecase 战绩查询与排行榜用例。
type StatsUsecase struct {
	repo StatsRepo
	log  *log.Helper
}

// NewStatsUsecase 创建 StatsUsecase。
func NewStatsUsecase(repo StatsRepo, logger log.Logger) *StatsUsecase {
	return &StatsUsecase{repo: repo, log: log.NewHelper(logger)}
}

// ListMatches 分页查询战绩。userID 为 nil 时查全部（当前无登录态）。
func (uc *StatsUsecase) ListMatches(ctx context.Context, page, pageSize int32, userID *uint64) ([]*MatchRecordView, int64, error) {
	uc.log.Infof("list matches, page=%d, pageSize=%d", page, pageSize)
	return uc.repo.ListMatchRecords(ctx, page, pageSize, userID)
}

// GetLeaderboard 获取排行榜。
func (uc *StatsUsecase) GetLeaderboard(ctx context.Context, limit int32) ([]*LeaderboardEntryView, error) {
	uc.log.Infof("get leaderboard, limit=%d", limit)
	return uc.repo.GetLeaderboard(ctx, limit)
}
