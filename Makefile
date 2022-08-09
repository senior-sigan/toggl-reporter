build_linux:
	GOOS=linux GOARCH=amd64 go build -o goreporter_web_LINUX64 goreporter/web

build:
	GOROOT=${GOROOT} GOPATH=${GOPATH} go build -o goreporter_web goreporter/web