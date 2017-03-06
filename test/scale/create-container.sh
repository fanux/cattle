docker -H :4001 run -d -l foo=bar -l app=whoami -l service=web -l type=online --name 0Whoami dev.reg.iflytek.com/devops/whoami:latest
docker -H :4001 run -d -l key=value -l app=fack -l service=compute -l type=offline --name 0Fack dev.reg.iflytek.com/devops/whoami:latest
