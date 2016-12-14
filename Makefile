coverage:
	@go test -coverprofile=/tmp/coverage && go tool cover -html=/tmp/coverage

test:
	@go test -v
