// Package service AI 内容加工坊（Core）HTTP 服务：实现 api/xsh/v1 ContentService。
// 模板与草稿 CRUD；GenerateCopy/GenerateCover/ProcessVideo 首期使用 Noop 时返回明确错误提示。

package service

import (
	"context"

	xshv1 "realworld/api/xsh/v1"
	"realworld/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// XshContentService 实现 xsh.v1.ContentService 的 HTTP 服务。
type XshContentService struct {
	uc  *biz.XshContentUsecase
	log *log.Helper
}

// NewXshContentService 创建内容加工坊 HTTP 服务实例。
func NewXshContentService(uc *biz.XshContentUsecase, logger log.Logger) *XshContentService {
	return &XshContentService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// ListTemplates 列表模板，可按 style_type 筛选。
func (s *XshContentService) ListTemplates(ctx context.Context, req *xshv1.ListTemplatesRequest) (*xshv1.ListTemplatesReply, error) {
	list, err := s.uc.ListTemplates(ctx, req.GetStyleType())
	if err != nil {
		s.log.Errorf("ListTemplates failed: %v", err)
		return nil, err
	}
	items := make([]*xshv1.Template, 0, len(list))
	for _, t := range list {
		items = append(items, &xshv1.Template{
			Id:              t.ID,
			Name:            t.Name,
			StyleType:       t.StyleType,
			ContentTemplate: t.ContentTemplate,
			SortOrder:       t.SortOrder,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
		})
	}
	return &xshv1.ListTemplatesReply{Items: items}, nil
}

// CreateTemplate 创建模板。
func (s *XshContentService) CreateTemplate(ctx context.Context, req *xshv1.CreateTemplateRequest) (*xshv1.CreateTemplateReply, error) {
	t := &biz.XshTemplate{
		Name:            req.GetName(),
		StyleType:       req.GetStyleType(),
		ContentTemplate: req.GetContentTemplate(),
		SortOrder:       req.GetSortOrder(),
	}
	created, err := s.uc.CreateTemplate(ctx, t)
	if err != nil {
		s.log.Errorf("CreateTemplate failed: %v", err)
		return nil, err
	}
	return &xshv1.CreateTemplateReply{Template: templateToProto(created)}, nil
}

// UpdateTemplate 更新模板；请求中非空字段覆盖已有值。
func (s *XshContentService) UpdateTemplate(ctx context.Context, req *xshv1.UpdateTemplateRequest) (*xshv1.UpdateTemplateReply, error) {
	existing, err := s.uc.GetTemplate(ctx, req.GetId())
	if err != nil || existing == nil {
		return nil, err
	}
	if req.GetName() != "" {
		existing.Name = req.GetName()
	}
	if req.GetStyleType() != "" {
		existing.StyleType = req.GetStyleType()
	}
	if req.GetContentTemplate() != "" {
		existing.ContentTemplate = req.GetContentTemplate()
	}
	// sort_order 允许为 0，始终按请求更新
	existing.SortOrder = req.GetSortOrder()
	if err := s.uc.UpdateTemplate(ctx, existing); err != nil {
		s.log.Errorf("UpdateTemplate failed: %v", err)
		return nil, err
	}
	return &xshv1.UpdateTemplateReply{Template: templateToProto(existing)}, nil
}

// DeleteTemplate 删除模板。
func (s *XshContentService) DeleteTemplate(ctx context.Context, req *xshv1.DeleteTemplateRequest) (*xshv1.DeleteTemplateReply, error) {
	if err := s.uc.DeleteTemplate(ctx, req.GetId()); err != nil {
		s.log.Errorf("DeleteTemplate failed: %v", err)
		return nil, err
	}
	return &xshv1.DeleteTemplateReply{Success: true}, nil
}

// CreateDraft 创建草稿，初始 status 为 draft。
func (s *XshContentService) CreateDraft(ctx context.Context, req *xshv1.CreateDraftRequest) (*xshv1.CreateDraftReply, error) {
	d := &biz.XshDraft{
		ProductID:  req.GetProductId(),
		TemplateID: req.GetTemplateId(),
		Status:     "draft",
	}
	created, err := s.uc.CreateDraft(ctx, d)
	if err != nil {
		s.log.Errorf("CreateDraft failed: %v", err)
		return nil, err
	}
	return &xshv1.CreateDraftReply{Draft: draftToProto(created)}, nil
}

// ListDrafts 分页列表草稿。
func (s *XshContentService) ListDrafts(ctx context.Context, req *xshv1.ListDraftsRequest) (*xshv1.ListDraftsReply, error) {
	page, pageSize := req.GetPage(), req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if page <= 0 {
		page = 1
	}
	items, total, err := s.uc.ListDrafts(ctx, page, pageSize, req.GetStatus(), req.GetProductId())
	if err != nil {
		s.log.Errorf("ListDrafts failed: %v", err)
		return nil, err
	}
	list := make([]*xshv1.Draft, 0, len(items))
	for _, d := range items {
		list = append(list, draftToProto(d))
	}
	return &xshv1.ListDraftsReply{Items: list, Total: total}, nil
}

// GetDraft 按 ID 查询草稿。
func (s *XshContentService) GetDraft(ctx context.Context, req *xshv1.GetDraftRequest) (*xshv1.GetDraftReply, error) {
	d, err := s.uc.GetDraft(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if d == nil {
		return &xshv1.GetDraftReply{}, nil
	}
	return &xshv1.GetDraftReply{Draft: draftToProto(d)}, nil
}

// UpdateDraft 更新草稿；请求中非空字段覆盖已有值。
func (s *XshContentService) UpdateDraft(ctx context.Context, req *xshv1.UpdateDraftRequest) (*xshv1.UpdateDraftReply, error) {
	existing, err := s.uc.GetDraft(ctx, req.GetId())
	if err != nil || existing == nil {
		return nil, err
	}
	if req.GetCopyText() != "" {
		existing.CopyText = req.GetCopyText()
	}
	if req.GetCoverImageUrl() != "" {
		existing.CoverImageURL = req.GetCoverImageUrl()
	}
	if req.GetVideoUrl() != "" {
		existing.VideoURL = req.GetVideoUrl()
	}
	if req.GetStatus() != "" {
		existing.Status = req.GetStatus()
	}
	if err := s.uc.UpdateDraft(ctx, existing); err != nil {
		s.log.Errorf("UpdateDraft failed: %v", err)
		return nil, err
	}
	return &xshv1.UpdateDraftReply{Draft: draftToProto(existing)}, nil
}

// DeleteDraft 删除草稿。
func (s *XshContentService) DeleteDraft(ctx context.Context, req *xshv1.DeleteDraftRequest) (*xshv1.DeleteDraftReply, error) {
	if err := s.uc.DeleteDraft(ctx, req.GetId()); err != nil {
		s.log.Errorf("DeleteDraft failed: %v", err)
		return nil, err
	}
	return &xshv1.DeleteDraftReply{Success: true}, nil
}

// GenerateCopy 调用 Gemini 生成文案；未配置时返回错误提示。
func (s *XshContentService) GenerateCopy(ctx context.Context, req *xshv1.GenerateCopyRequest) (*xshv1.GenerateCopyReply, error) {
	copyText, err := s.uc.GenerateCopy(ctx, req.GetDraftId(), req.GetStyle())
	if err != nil {
		s.log.Warnf("GenerateCopy failed: %v", err)
		return nil, err
	}
	return &xshv1.GenerateCopyReply{CopyText: copyText}, nil
}

// GenerateCover 合成封面；未配置时返回错误提示。
func (s *XshContentService) GenerateCover(ctx context.Context, req *xshv1.GenerateCoverRequest) (*xshv1.GenerateCoverReply, error) {
	coverURL, err := s.uc.GenerateCover(ctx, req.GetDraftId())
	if err != nil {
		s.log.Warnf("GenerateCover failed: %v", err)
		return nil, err
	}
	return &xshv1.GenerateCoverReply{CoverImageUrl: coverURL}, nil
}

// ProcessVideo 视频去重；未配置时返回错误提示。
func (s *XshContentService) ProcessVideo(ctx context.Context, req *xshv1.ProcessVideoRequest) (*xshv1.ProcessVideoReply, error) {
	draft, err := s.uc.GetDraft(ctx, req.GetDraftId())
	if err != nil || draft == nil {
		return nil, err
	}
	videoURL, err := s.uc.ProcessVideo(ctx, req.GetDraftId(), draft.VideoURL, req.GetOptions())
	if err != nil {
		s.log.Warnf("ProcessVideo failed: %v", err)
		return nil, err
	}
	return &xshv1.ProcessVideoReply{VideoUrl: videoURL}, nil
}

func templateToProto(t *biz.XshTemplate) *xshv1.Template {
	if t == nil {
		return nil
	}
	return &xshv1.Template{
		Id:              t.ID,
		Name:            t.Name,
		StyleType:       t.StyleType,
		ContentTemplate: t.ContentTemplate,
		SortOrder:       t.SortOrder,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

func draftToProto(d *biz.XshDraft) *xshv1.Draft {
	if d == nil {
		return nil
	}
	return &xshv1.Draft{
		Id:              d.ID,
		ProductId:       d.ProductID,
		TemplateId:      d.TemplateID,
		CopyText:        d.CopyText,
		CoverImageUrl:   d.CoverImageURL,
		VideoUrl:        d.VideoURL,
		Status:          d.Status,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}
