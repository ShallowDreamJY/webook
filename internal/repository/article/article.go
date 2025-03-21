package article

import (
	"context"
	"gorm.io/gorm"
	"webook/internal/domain"
	dao "webook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// 存储同步
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncV1(ctx context.Context, art domain.Article) (int64, error)
	SyncV2(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDao

	// v1
	readerDAO dao.ReaderDAO
	authorDAO dao.AuthorDAO

	db *gorm.DB
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

// 确保保存到线上库和制作库同时成功
// 开启数据库事务
func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	// 利用tx来构建
	author := dao.NewAuthorDAO(tx)
	reader := dao.NewReaderDAO(tx)

	var (
		id  = art.Id
		err error
	)
	artEn := c.toEntity(art)
	// 保存制作库
	if art.Id == 0 {
		id, err = author.Insert(ctx, artEn)
	} else {
		err = author.UpdateById(ctx, artEn)
	}
	if err != nil {
		return 0, err
	}
	// 保存线上库
	err = reader.Upsert(ctx, artEn)
	tx.Commit()
	return id, err
}

func NewArticleRepository(dao dao.ArticleDao,
	readerDAO dao.ReaderDAO,
	authorDAO dao.AuthorDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao:       dao,
		readerDAO: readerDAO,
		authorDAO: authorDAO}
}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artEn := c.toEntity(art)
	// 保存制作库
	if art.Id == 0 {
		id, err = c.authorDAO.Insert(ctx, artEn)
	} else {
		err = c.authorDAO.UpdateById(ctx, artEn)
	}
	if err != nil {
		return 0, err
	}

	// 保存线上库
	err = c.readerDAO.Upsert(ctx, artEn)
	return id, err
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
