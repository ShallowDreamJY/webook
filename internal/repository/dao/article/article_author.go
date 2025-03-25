package article

import (
	"context"
	"gorm.io/gorm"
)

type AuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}

type GORMAuthorDAO struct {
	db *gorm.DB
}

func (G GORMAuthorDAO) Insert(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (G GORMAuthorDAO) UpdateById(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func NewAuthorDAO(db *gorm.DB) AuthorDAO {
	return &GORMAuthorDAO{
		db: db,
	}
}
