FROM 192.168.86.106/devops/golang:1.7-godep
COPY cattle .
CMD ./cattle --help
