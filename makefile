GO_LDFLAGS=-ldflags " -w"

TAG=latest
PREFIX=dhub.yunpro.cn/barnett/alert

build: ## build the go packages
	@echo "🐳 $@"
	@go build -a -installsuffix cgo ${GO_LDFLAGS} .

image: clean
	@echo "🐳 $@"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ${GO_LDFLAGS} .
	@docker build -t $(PREFIX):$(TAG) .
	@docker push $(PREFIX):$(TAG)
	
clean:
	@rm -f alertCenter