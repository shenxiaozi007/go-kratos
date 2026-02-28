package article

import (
	"time"

	"gorm.io/gorm"
)

type ArticleStatus uint8

// 0=草稿 1=已发布
const (
	Draft ArticleStatus = iota
	Published
)

type Article struct {
	gorm.Model
	id         int64         `gorm:"primary_key"`
	title      string        `gorm:"type:varchar(100);not null"`
	content    string        `gorm:"type:text;not null"`
	status     ArticleStatus `gorm:"type:tinyint(1);not null"`
	categoryID int64
	viewCount  int
	likeCount  int
}

func NewArticle(title, content string, categoryID int64) (*Article, error) {
	if title == "" {
		return nil, ErrTitleNotFound
	}

	if content == "" {
		return nil, ErrContentNotFound
	}

	return &Article{
		title:      title,
		content:    content,
		status:     Draft,
		categoryID: categoryID,
		viewCount:  0,
		likeCount:  0,
	}, nil
}

// 发布文章
func (a *Article) Publish() error {
	if a.status == Published {
		return ErrAlreadyPub
	}
	a.status = Published
	a.UpdatedAt = time.Now()
	return nil
}

// 阅读文章
func (a *Article) View() {
	a.viewCount++
}

func (a *Article) Like() {
	a.likeCount++
}

// 暴露setter
func (a *Article) ID() int64 {
	return a.id
}

func (a *Article) Title() string {
	return a.title
}

func (a *Article) Content() string {
	return a.content
}

func (a *Article) Status() ArticleStatus {
	return a.status
}

func (a *Article) ViewCount() int {
	return a.viewCount
}

func (a *Article) LikeCount() int {
	return a.likeCount
}

func (a *Article) CategoryID() int64 {
	return a.categoryID
}
