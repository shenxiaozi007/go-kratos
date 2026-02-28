package service

import (
	"context"
	"realworld/internal/biz"
	"realworld/internal/domian/article"

	pb "realworld/api/blog/v1"
)

type ArticleService struct {
	pb.UnimplementedArticleServer
	uc *biz.ArticleUseCase
}

func NewArticleService(uc *biz.ArticleUseCase) *ArticleService {
	return &ArticleService{
		uc: uc,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.ArticleReply, error) {
	ar, err := s.uc.CreateArticle(ctx, req.Title, req.Content, req.CategoryId)
	if err != nil {
		return nil, err
	}
	return toArticleReply(ar), nil
}

// 获取文章的时候。顺便加view
func (s *ArticleService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.ArticleReply, error) {
	ar, err := s.uc.ViewArticle(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toArticleReply(ar), nil
}

func (s *ArticleService) ListArticle(ctx context.Context, req *pb.ListArticleRequest) (*pb.ListArticleReply, error) {

	ars, err := s.uc.ListPublishedArticles(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, err
	}

	//var articleReplyList []*pb.ArticleReply
	//for _, item := range ars {
	//	articleReplyList = append(articleReplyList, toArticleReply(item))
	//}
	//
	//return &pb.ListArticleReply{
	//	List: articleReplyList,
	//}, nil
	// 另一种写法
	resp := &pb.ListArticleReply{}

	for _, item := range ars {
		resp.List = append(resp.List, toArticleReply(item))
	}
	return resp, nil

}

// 发布文章
func (s *ArticleService) PublishArticle(ctx context.Context, req *pb.PublishArticleRequest) (*pb.EmptyReply, error) {
	if err := s.uc.PublishArticle(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

// dto转换
func toArticleReply(a *article.Article) *pb.ArticleReply {
	return &pb.ArticleReply{
		Title:     a.Title(),
		Content:   a.Content(),
		Status:    int64(a.Status()),
		ViewCount: int64(a.ViewCount()),
		LikeCount: int64(a.LikeCount()),
	}
}
