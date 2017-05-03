sh destroy-all-containers.sh

echo "create container"
docker -H :4001 run -d -e constraint:pool==true -e affinity:app!=online -e afinity:app!=offline -l app=online --name iat dev.reg.iflytek.com/devops/whoami:latest

../../swarm scale -e constraint:pool==true -f app==online -n 8

echo "show containers"
sh show-containers.sh
