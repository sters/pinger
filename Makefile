GO_ENV := GO111MODULE=on CGO_ENABLED=0

.PHONY: init tidy test
init: 
	@${GO_ENV} go mod init
tidy: 
	@${GO_ENV} go mod tidy
test: 
	@${GO_ENV} CGO_ENABLED=1 go test -v -race -cover ./...