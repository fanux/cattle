docker -H :4001 run -d -e affinity:app!=online -e afinity:app!=offline -l app=online --name iat dev.reg.iflytek.com/devops/whoami:latest
docker -H :4001 run -d -e affinity:app!=online -e affinity:app!=offline -l app=offline --name nodemanager dev.reg.iflytek.com/devops/whoami:latest

../../swarm scale -f app==online -n 3
../../swarm scale -f app==offline -n 3

echo "show containers:"
sh ./show-containers.sh
