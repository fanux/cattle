echo "clean all containers"
sh ../destroy-all-containers.sh

echo "create containers"
docker -H :4001 run -d -e PRIORITY=1 -e affinity:type!=offline  -l type=online --name iat dev.reg.iflytek.com/devops/whoami:latest
docker -H :4001 run -d  -e PRIORITY=9 -e affinity:type!=online -e affinity:type!=offline -l type=offline --name nodemanager dev.reg.iflytek.com/devops/whoami:latest
../../../swarm scale -f type==offline -n 6

echo "show containers:"
docker -H :4001 ps -a
echo "seize low priority resource, applots == 4"
../../../swarm scale -e constraint:GPU==true -e applots=4 -e affinity:type!=offline -f type==online -n 3

echo "show containers"
docker -H :4001 ps -a
