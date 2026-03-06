// Package biz AI 内容加工坊（Core）业务层：模板、草稿及文案/封面/视频生成的领域实体与用例。
// 文案生成由外部 Gemini 3 Flash 驱动，本层仅定义接口与编排。

package biz

import (
	"context"
	"errors"
)

// 内容加工坊可选能力未注入时的错误，便于上层返回明确提示。
var (
	ErrCopyGeneratorNotAvailable   = errors.New("文案生成服务未配置")
	ErrCoverGeneratorNotAvailable  = errors.New("封面合成服务未配置")
	ErrVideoProcessorNotAvailable  = errors.New("视频去重服务未配置")
)

// XshTemplate 文案模板领域实体，对应表 xsh_content_templates。
type XshTemplate struct {
	ID              uint64
	Name            string
	StyleType       string
	ContentTemplate string
	SortOrder       int32
	CreatedAt       int64
	UpdatedAt       int64
}

// XshDraft 内容草稿领域实体，对应表 xsh_content_drafts。
// Status: draft | ready | published。
type XshDraft struct {
	ID               uint64
	ProductID        uint64
	TemplateID       uint64
	CopyText         string
	CoverImageURL    string
	VideoURL         string
	Status           string
	CreatedAt        int64
	UpdatedAt        int64
}

// XshTemplateRepo 模板仓储接口。
type XshTemplateRepo interface {
	List(ctx context.Context, styleType string) ([]*XshTemplate, error)
	GetByID(ctx context.Context, id uint64) (*XshTemplate, error)
	Create(ctx context.Context, t *XshTemplate) (*XshTemplate, error)
	Update(ctx context.Context, t *XshTemplate) error
	Delete(ctx context.Context, id uint64) error
}

// XshDraftRepo 草稿仓储接口。
type XshDraftRepo interface {
	GetByID(ctx context.Context, id uint64) (*XshDraft, error)
	List(ctx context.Context, page, pageSize int32, status string, productID uint64) ([]*XshDraft, int64, error)
	Create(ctx context.Context, d *XshDraft) (*XshDraft, error)
	Update(ctx context.Context, d *XshDraft) error
	Delete(ctx context.Context, id uint64) error
}

// CopyGenerator 文案生成器接口，由 content 用例调用；实际实现可调用 Gemini 3 Flash。
type CopyGenerator interface {
	GenerateCopy(ctx context.Context, draftID uint64, style string, productTitle string, templateContent string) (string, error)
}

// CoverGenerator 封面合成器接口。
type CoverGenerator interface {
	GenerateCover(ctx context.Context, draftID uint64, productImageURL string) (string, error)
}

// VideoProcessor 视频去重处理接口。
type VideoProcessor interface {
	ProcessVideo(ctx context.Context, draftID uint64, videoURL string, options string) (string, error)
}

// XshContentUsecase AI 内容加工坊用例。
type XshContentUsecase struct {
	templateRepo   XshTemplateRepo
	draftRepo      XshDraftRepo
	productRepo    XshProductRepo
	copyGen        CopyGenerator
	coverGen       CoverGenerator
	videoProcessor VideoProcessor
}

// NewXshContentUsecase 创建内容加工坊用例。copyGen/coverGen/videoProcessor 可为 nil，对应能力不可用时返回明确错误。
func NewXshContentUsecase(
	templateRepo XshTemplateRepo,
	draftRepo XshDraftRepo,
	productRepo XshProductRepo,
	copyGen CopyGenerator,
	coverGen CoverGenerator,
	videoProcessor VideoProcessor,
) *XshContentUsecase {
	return &XshContentUsecase{
		templateRepo:   templateRepo,
		draftRepo:      draftRepo,
		productRepo:    productRepo,
		copyGen:        copyGen,
		coverGen:       coverGen,
		videoProcessor: videoProcessor,
	}
}

// ListTemplates 列表模板，可按 style_type 筛选。
func (uc *XshContentUsecase) ListTemplates(ctx context.Context, styleType string) ([]*XshTemplate, error) {
	return uc.templateRepo.List(ctx, styleType)
}

// GetTemplate 按 ID 查询模板。
func (uc *XshContentUsecase) GetTemplate(ctx context.Context, id uint64) (*XshTemplate, error) {
	return uc.templateRepo.GetByID(ctx, id)
}

// CreateTemplate 创建模板。
func (uc *XshContentUsecase) CreateTemplate(ctx context.Context, t *XshTemplate) (*XshTemplate, error) {
	return uc.templateRepo.Create(ctx, t)
}

// UpdateTemplate 更新模板。
func (uc *XshContentUsecase) UpdateTemplate(ctx context.Context, t *XshTemplate) error {
	return uc.templateRepo.Update(ctx, t)
}

// DeleteTemplate 删除模板。
func (uc *XshContentUsecase) DeleteTemplate(ctx context.Context, id uint64) error {
	return uc.templateRepo.Delete(ctx, id)
}

// CreateDraft 创建草稿，初始 status 为 draft。
func (uc *XshContentUsecase) CreateDraft(ctx context.Context, d *XshDraft) (*XshDraft, error) {
	return uc.draftRepo.Create(ctx, d)
}

// ListDrafts 分页列表草稿，支持 status、product_id 筛选。
func (uc *XshContentUsecase) ListDrafts(ctx context.Context, page, pageSize int32, status string, productID uint64) ([]*XshDraft, int64, error) {
	return uc.draftRepo.List(ctx, page, pageSize, status, productID)
}

// GetDraft 按 ID 查询草稿。
func (uc *XshContentUsecase) GetDraft(ctx context.Context, id uint64) (*XshDraft, error) {
	return uc.draftRepo.GetByID(ctx, id)
}

// UpdateDraft 更新草稿（文案、封面、视频、status 等）。
func (uc *XshContentUsecase) UpdateDraft(ctx context.Context, d *XshDraft) error {
	return uc.draftRepo.Update(ctx, d)
}

// DeleteDraft 删除草稿。
func (uc *XshContentUsecase) DeleteDraft(ctx context.Context, id uint64) error {
	return uc.draftRepo.Delete(ctx, id)
}

// GenerateCopy 调用 CopyGenerator（如 Gemini）生成文案；若未注入则返回 ErrCopyGeneratorNotAvailable。
func (uc *XshContentUsecase) GenerateCopy(ctx context.Context, draftID uint64, style string) (string, error) {
	if uc.copyGen == nil {
		return "", ErrCopyGeneratorNotAvailable
	}
	draft, err := uc.draftRepo.GetByID(ctx, draftID)
	if err != nil || draft == nil {
		return "", err
	}
	product, _ := uc.productRepo.GetByID(ctx, draft.ProductID)
	productTitle := ""
	if product != nil {
		productTitle = product.Title
	}
	templateContent := ""
	if draft.TemplateID > 0 {
		if t, _ := uc.templateRepo.GetByID(ctx, draft.TemplateID); t != nil {
			templateContent = t.ContentTemplate
		}
	}
	return uc.copyGen.GenerateCopy(ctx, draftID, style, productTitle, templateContent)
}

// GenerateCover 调用 CoverGenerator 合成封面；若未注入则返回 ErrCoverGeneratorNotAvailable。
func (uc *XshContentUsecase) GenerateCover(ctx context.Context, draftID uint64) (string, error) {
	if uc.coverGen == nil {
		return "", ErrCoverGeneratorNotAvailable
	}
	draft, err := uc.draftRepo.GetByID(ctx, draftID)
	if err != nil || draft == nil {
		return "", err
	}
	product, _ := uc.productRepo.GetByID(ctx, draft.ProductID)
	productImageURL := ""
	if product != nil {
		productImageURL = product.ImageURL
	}
	return uc.coverGen.GenerateCover(ctx, draftID, productImageURL)
}

// ProcessVideo 调用 VideoProcessor 做视频去重；若未注入则返回 ErrVideoProcessorNotAvailable。
func (uc *XshContentUsecase) ProcessVideo(ctx context.Context, draftID uint64, videoURL string, options string) (string, error) {
	if uc.videoProcessor == nil {
		return "", ErrVideoProcessorNotAvailable
	}
	return uc.videoProcessor.ProcessVideo(ctx, draftID, videoURL, options)
}
