package service

import (
	"context"

	statsv1 "realworld/api/stats/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// StatsService 实现 stats.v1.StatsServiceServer，提供历史战绩与排行榜接口。
type StatsService struct {
	statsv1.UnimplementedStatsServiceServer

	uc  *biz.StatsUsecase
	log *log.Helper
}

// NewStatsService 创建 StatsService。
func NewStatsService(uc *biz.StatsUsecase, logger log.Logger) *StatsService {
	return &StatsService{uc: uc, log: log.NewHelper(logger)}
}

// ListMatches 分页查询历史战绩。当前无登录态，查询全部记录。
func (s *StatsService) ListMatches(ctx context.Context, req *statsv1.ListMatchesRequest) (*statsv1.ListMatchesReply, error) {
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	records, total, err := s.uc.ListMatches(ctx, page, pageSize, nil)
	if err != nil {
		s.log.Errorf("ListMatches failed: %v", err)
		return nil, err
	}
	reply := &statsv1.ListMatchesReply{
		Records:  make([]*statsv1.MatchRecord, 0, len(records)),
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	for _, r := range records {
		reply.Records = append(reply.Records, &statsv1.MatchRecord{
			SessionId:       r.SessionID,
			UserId:          int64(r.UserID),
			Mode:            statsv1.GameMode(r.Mode),
			MapId:           r.MapID,
			IsWin:           r.IsWin,
			Player1Kills:    r.Player1Kills,
			Player2Kills:    r.Player2Kills,
			EnemyKills:      r.EnemyKills,
			DurationSeconds: r.DurationSecond,
			CreatedAt:       r.CreatedAt,
		})
	}
	return reply, nil
}

// GetLeaderboard 返回排行榜。
func (s *StatsService) GetLeaderboard(ctx context.Context, req *statsv1.GetLeaderboardRequest) (*statsv1.GetLeaderboardReply, error) {
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 50
	}
	entries, err := s.uc.GetLeaderboard(ctx, limit)
	if err != nil {
		s.log.Errorf("GetLeaderboard failed: %v", err)
		return nil, err
	}
	reply := &statsv1.GetLeaderboardReply{
		Entries: make([]*statsv1.LeaderboardEntry, 0, len(entries)),
	}
	for _, e := range entries {
		reply.Entries = append(reply.Entries, &statsv1.LeaderboardEntry{
			UserId:       int64(e.UserID),
			Username:     e.Username,
			TotalMatches: int32(e.TotalMatches),
			TotalWins:    int32(e.TotalWins),
			WinRate:      e.WinRate,
		})
	}
	return reply, nil
}
