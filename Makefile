HONY: mock
ck:
	@mockgen -source="internal/repository/user.go" -package=repomocks -destination="internal/repository/mocks/user.mock.go"
	@mockgen -source="internal/service/user.go" -package=svcmocks -destination="internal/service/mocks/user.mock.go"
	@mockgen -source="internal/service/article.go" -package=svcmocks -destination="internal/service/mocks/article.mock.go"
	@mockgen -source="internal/repository/code.go" -package=repomocks -destination="internal/repository/mocks/code.mock.go"

	@mockgen -source="internal/repository/article/article_author.go" -package=artrepomocks -destination="internal/repository/article/mocks/article_author.mock.go"
	@mockgen -source="internal/repository/article/article_reader.go" -package=artrepomocks -destination="internal/repository/article/mocks/article_reader.mock.go"
	@mockgen -source="internal/repository/article/article.go" -package=artrepomocks -destination="internal/repository/article/mocks/article.mock.go"
	@mockgen -source="internal/service/article.go" -package=svcmocks -destination="internal/service/mocks/article.mock.go"
	@go mod tidy