FROM 192.168.86.106/devops/alpine:3.4
COPY /drone/cattle /bin
CMD cattle --help
