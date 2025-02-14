HONY: mock
ck:
	@mockgen -source="internal/repository/user.go" -package=repomocks -destination="internal/repository/mocks/user.mock.go"
	@mockgen -source="internal/repository/code.go" -package=repomocks -destination="internal/repository/mocks/code.mock.go"
	@go mod tidy