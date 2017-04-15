GO_LDFLAGS=-ldflags " -w"

TAG=0415
# PREFIX=barnettzqg/alert-center
PREFIX = goodrain.me/8439cf79b5c6_barnettZQG_alertCenter
build: ## build the go packages
	@echo "🐳 $@"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ${GO_LDFLAGS} .

image: clean
	@echo "🐳 $@"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ${GO_LDFLAGS} .
	@docker build -t $(PREFIX):$(TAG) .
	@docker push $(PREFIX):$(TAG)
	@rm -f alertCenter
clean:
	@rm -f alertCenter