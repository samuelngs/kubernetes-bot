
all:
	CGO_ENABLED=0 go build -a -v -installsuffix cgo -o bin/kubebot cmd/main.go
