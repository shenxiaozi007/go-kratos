package article

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, article *Article) error
	Update(ctx context.Context, article *Article) error
	FindByID(ctx context.Context, id int64) (*Article, error)
	ListPublished(ctx context.Context, offset, limit int) ([]*Article, error)
}

func RebuildArticle(
	id int64,
	title string,
	content string,
	categoryID int64,
	status ArticleStatus,
	viewCount int,
	likeCount int,
) *Article {
	return &Article{
		id:         id,
		title:      title,
		content:    content,
		categoryID: categoryID,
		status:     status,
		viewCount:  viewCount,
		likeCount:  likeCount,
	}
}
