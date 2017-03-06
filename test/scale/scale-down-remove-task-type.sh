echo "\n\n"
echo "Show containers:"
sh show-containers.sh

echo "\n\n"
echo "Destroy all containers"
sh destroy-all-containers.sh
sh show-containers.sh

echo "\n\n"
echo "Create test containers"
sh create-container.sh
sh show-containers.sh

echo "\n\n"
echo "Scale up containers"
sh scale-up-foo-bar.sh
sh show-containers.sh

echo "\n\n"
echo "Scale down containers"
../../swarm scale -f foo==bar -n -2
../../swarm scale -f foo==bar -n -3
sh show-containers.sh
