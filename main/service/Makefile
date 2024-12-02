
NAME=ss-proxy
BINDIR=bin

linux-amd64:
	GOARCH=amd64 GOOS=linux  CGO_ENABLED=0 go build -o deploy/clash-ws
clean:
	rm -rf $(BINDIR)/*
