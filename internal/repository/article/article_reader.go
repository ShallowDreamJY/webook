package article

import (
	"context"
	"webook/internal/domain"
)

type ArticleReaderRepository interface {
	// 有就更新，无则新建
	Save(ctx context.Context, art domain.Article) (int64, error)
}
