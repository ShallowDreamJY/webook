package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUnit8() uint8 {
	return uint8(s)
}

func (s ArticleStatus) NonPublish() bool {
	return s != ArticleStatusPublished
}

func (s ArticleStatus) Valid() bool {
	return s.ToUnit8() > 0
}

func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusPrivate:
		return "Private"
	case ArticleStatusUnpublished:
		return "Unpublished"
	case ArticleStatusPublished:
		return "Published"
	default:
		return "Unknown"
	}
}

type Author struct {
	Id   int64
	Name string
}
