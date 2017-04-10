echo "clean all containers"
sh ../destroy-all-containers.sh

echo "create container"
docker -H :4001 run  -e STOP_TIMEOUT=10 -l key=value dev.reg.iflytek.com/test/signal:latest
