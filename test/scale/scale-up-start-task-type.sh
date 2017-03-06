sh scale-down-foo-bar.sh
echo "\n\n"

../../swarm scale -f foo==bar -e TASK_TYPE=start -n 3
../../swarm scale -f foo==bar -e TASK_TYPE=start -n 2
echo "\n\n"

sh show-containers.sh
