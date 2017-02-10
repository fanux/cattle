FROM dev.reg.iflytek.com/devops/golang:1.7-alpine
COPY cattle $GOPATH/bin
CMD cattle --help
