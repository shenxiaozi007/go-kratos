// Package biz 内容加工坊可选能力占位实现，用于 Wire 注入。
// 当未配置 Gemini/封面合成/视频去重时，使用本占位实现，返回明确错误提示。

package biz

import "context"

// NoopCopyGenerator 占位文案生成器，返回 ErrCopyGeneratorNotAvailable。
type NoopCopyGenerator struct{}

// GenerateCopy 未配置时返回错误。
func (NoopCopyGenerator) GenerateCopy(ctx context.Context, draftID uint64, style, productTitle, templateContent string) (string, error) {
	return "", ErrCopyGeneratorNotAvailable
}

// NoopCoverGenerator 占位封面合成器。
type NoopCoverGenerator struct{}

// GenerateCover 未配置时返回错误。
func (NoopCoverGenerator) GenerateCover(ctx context.Context, draftID uint64, productImageURL string) (string, error) {
	return "", ErrCoverGeneratorNotAvailable
}

// NoopVideoProcessor 占位视频处理器。
type NoopVideoProcessor struct{}

// ProcessVideo 未配置时返回错误。
func (NoopVideoProcessor) ProcessVideo(ctx context.Context, draftID uint64, videoURL, options string) (string, error) {
	return "", ErrVideoProcessorNotAvailable
}

// NewNoopCopyGenerator 返回占位文案生成器，供 Wire 注入。
func NewNoopCopyGenerator() CopyGenerator {
	return &NoopCopyGenerator{}
}

// NewNoopCoverGenerator 返回占位封面合成器，供 Wire 注入。
func NewNoopCoverGenerator() CoverGenerator {
	return &NoopCoverGenerator{}
}

// NewNoopVideoProcessor 返回占位视频处理器，供 Wire 注入。
func NewNoopVideoProcessor() VideoProcessor {
	return &NoopVideoProcessor{}
}
