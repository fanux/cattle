FROM 10.1.86.51/devops/golang:1.7-alpine
COPY cattle $GOPATH/bin
RUN touch /a
RUN mv /a /a
CMD cattle --help
