package biz

import (
	"context"
	"realworld/internal/domian/article"
)

type ArticleUseCase struct {
	repo article.Repository
}

func NewArticleUseCase(repo article.Repository) *ArticleUseCase {
	return &ArticleUseCase{repo: repo}
}

func (uc *ArticleUseCase) CreateArticle(ctx context.Context, title, content string, categoryID int64) (*article.Article, error) {

	ar, err := article.NewArticle(title, content, categoryID)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Save(ctx, ar); err != nil {
		return nil, err
	}
	return ar, nil
}

// 发布文章
func (uc *ArticleUseCase) PublishArticle(ctx context.Context, id int64) error {
	ar, err := uc.repo.FindByID(ctx, id)

	if err != nil {
		return err
	}
	if err := ar.Publish(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, ar)
}

// 阅读文章
func (uc *ArticleUseCase) ViewArticle(ctx context.Context, id int64) (*article.Article, error) {
	ar, err := uc.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	ar.View()

	if err := uc.repo.Update(ctx, ar); err != nil {
		return nil, err
	}
	return ar, nil
}

// 点赞文章
func (uc *ArticleUseCase) LikeArticle(ctx context.Context, id int64) error {
	a, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	a.Like()

	return uc.repo.Update(ctx, a)
}

// 获取已经发布的文章
func (uc *ArticleUseCase) ListPublishedArticles(ctx context.Context, offset, limit int) ([]*article.Article, error) {
	return uc.repo.ListPublished(ctx, offset, limit)
}
