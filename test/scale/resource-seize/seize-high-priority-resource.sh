sh create-resource-seize-container.sh

echo "seize low priority resource"
../../../swarm scale -e constraint:GPU==true -e applots=1 -e affinity:type!=offline -f type==online -n 3

echo "show containers"
docker -H :4001 ps -a

echo "seize high priority resource"
../../../swarm scale -e constraint:GPU==true -e applots=1 -e affinity:type!=online -f type==offline -n 3

echo "show containers"
docker -H :4001 ps -a
