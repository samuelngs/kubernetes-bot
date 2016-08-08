
src:
	@CGO_ENABLED=0 go build -a -v -installsuffix cgo -o bin/kubebot cmd/main.go

docker:
	@docker build --no-cache --tag kubebot:latest --file docker/Dockerfile .

all: src docker

.PHONY: src docker all
