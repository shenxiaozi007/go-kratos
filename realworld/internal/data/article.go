package data

import (
	"context"
	"realworld/internal/domian/article"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ArticlePO struct {
	// 主键
	ID int64 `gorm:"primaryKey;autoIncrement"`
	// 标题
	Title string `gorm:"type:varchar(100);not null"`
	// 内容
	Content string `gorm:"type:text;not null"`
	// 状态
	Status uint8
	// 分类ID
	CategoryID int64
	// 阅读数
	ViewCount int
	// 点赞数
	LikeCount int
}

func (ArticlePO) TableName() string {
	return "article"
}

type articleRepo struct {
	db *Data
}

func NewArticleRepo(db *gorm.DB, logger log.Logger) article.Repository {
	return &articleRepo{
		db: &Data{
			DB: db,
		},
	}
}

// domain to Po
func toArticlePO(a *article.Article) *ArticlePO {
	return &ArticlePO{
		ID:         a.ID(),
		Title:      a.Title(),
		Content:    a.Content(),
		Status:     uint8(a.Status()),
		CategoryID: a.CategoryID(),
		ViewCount:  a.ViewCount(),
		LikeCount:  a.LikeCount(),
	}
}

func toArticleDomain(po *ArticlePO) *article.Article {
	return article.RebuildArticle(
		po.ID,
		po.Title,
		po.Content,
		po.CategoryID,
		article.ArticleStatus(po.Status),
		po.ViewCount,
		po.LikeCount,
	)
}

func (a *articleRepo) Save(ctx context.Context, ar *article.Article) error {
	//
	po := toArticlePO(ar)

	return a.db.DB.WithContext(ctx).Create(po).Error
}

func (a *articleRepo) Update(ctx context.Context, ar *article.Article) error {
	po := toArticlePO(ar)

	return a.db.DB.WithContext(ctx).Model(&ArticlePO{}).Where("id = ?", po.ID).Updates(po).Error
}

func (a *articleRepo) FindByID(ctx context.Context, id int64) (*article.Article, error) {
	var pos ArticlePO
	if err := a.db.DB.WithContext(ctx).First(&pos, id).Error; err != nil {
		return nil, err
	}

	return toArticleDomain(&pos), nil
}

func (a *articleRepo) ListPublished(ctx context.Context, offset, limit int) ([]*article.Article, error) {
	var pos []ArticlePO
	err := a.db.DB.WithContext(ctx).Where("status = ?", article.Published).Order("id desc").Offset(offset).Limit(limit).Find(&pos).Error
	if err != nil {
		return nil, err
	}

	articles := make([]*article.Article, 0, len(pos))

	for _, po := range pos {
		articles = append(articles, toArticleDomain(&po))
	}
	return articles, nil
}
