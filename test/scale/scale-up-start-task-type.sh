sh scale-down-stop-task-type.sh
echo "\n\n"

../../swarm scale -f foo==bar -e TASK_TYPE=start -n 3
../../swarm scale -f foo==bar -e TASK_TYPE=start -n 1
echo "\n\n"
sh show-containers.sh


echo "\n\n"
../../swarm scale -f foo==bar -e TASK_TYPE=start -n 3
echo "\n\n"
sh show-containers.sh
