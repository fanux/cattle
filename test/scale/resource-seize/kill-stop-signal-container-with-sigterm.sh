echo "show containers"
docker -H :4001 ps -a

echo "scale down to 0"
../../../swarm scale  -f key==value -n -1
