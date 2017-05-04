FROM dev.reg.iflytek.com/base/golang:1.8.0
COPY main.go .
RUN go build -o signal && cp signal $GOPATH/bin
CMD ["signal"] 
#CMD signal    这种会以shell的方式启动,sh是1号进程,上面那样样直接exec
