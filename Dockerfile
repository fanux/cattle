FROM 10.1.86.51/devops/golang:1.7-alpine
COPY cattle $GOPATH/bin
RUN mkdir /confd && touch /confd/confd && echo "asdfas">> /confd/confd
RUN mv /confd/confd /confd/confd
CMD cattle --help
