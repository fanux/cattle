FROM 192.168.86.106/devops/golang:1.7-alpine
COPY cattle $GOPATH/bin
CMD cattle --help
