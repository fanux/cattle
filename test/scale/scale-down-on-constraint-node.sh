sh destroy-all-containers.sh
sh init-iat-containers.sh

echo "scale down on constraint node"
../../swarm scale -f app==online -e constraint:pool==true -n -4
../../swarm scale -f app==offline -e constraint:pool==true -n -4

echo "show containers"
sh show-containers.sh
