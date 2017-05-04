sh   destroy-all-containers.sh

docker -H :4001 run -d -e constraint:host==10.1.86.201 -e affinity:app!=online -e affinity:app!=offline -l app=online --name iat dev.reg.iflytek.com/devops/whoami:latest
docker -H :4001 run -d -e constraint:host==10.1.86.208 -e affinity:app!=online -e affinity:app!=offline -l app=offline --name nodemanager dev.reg.iflytek.com/devops/whoami:latest

../../swarm scale -e constraint:pool==true -f app==online -n 3
../../swarm scale  -e constraint:pool==true -f app==offline -n 3

echo "show containers:"
sh ./show-containers.sh
