package article

import (
	"context"
	"gorm.io/gorm"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, art Article) error
}

type GORMReaderDAO struct {
	db *gorm.DB
}

func (G GORMReaderDAO) Upsert(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

type PublishArticle struct {
	Article
}

func NewReaderDAO(db *gorm.DB) ReaderDAO {
	return &GORMReaderDAO{db: db}
}
