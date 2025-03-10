package integration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)
// 测试套件
type ArticleTestSuite struct {
	suite.Suite
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string

		// 集成测试准备数据
		before func(t *testing.T)
		// 集成测试验证数据
		after func(t *testing.T)
		// 预期输入
		req Article

		// HTTP 响应码
		wantCode int
		// HTTP响应中带上文章 ID
		wantRes Result[int64]
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("hello test suit")
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	title   string `json:"title"`
	content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
