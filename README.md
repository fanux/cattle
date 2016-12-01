## Why cattle
1. Swarm not support scale command or API.
2. Scale up is easy, but when some high priority service want to seize low priority resources, how to decide witch  service to scale down?

Cattle sovle those problems.

Why not swarmkit or swarmmod? 

1. Not support `--net=host` (maybe support later)
2. Current not support some filter (later will support)
3. Swarmkit or swarmmod add more api to docker engine, witch I don't want to use.
4. Most importent is, swarm is simple.
5. We don't need more concept.

## Container priority, min number

```
$ docker run -e PRIORITY=10 -e MIN_NUMBER=3 -l service=online --name nginx nginx:latest
$ docker run -e PRIORITY=1 -e MIN_NUMBER=1 -l service=offline --name nginx httpd:latest
```

## Scale up or down
Cli

Scale up: suggest use docker compose service as a scale unit.
```
$ cattle scale [[env|label] cotainer number]
```
```
$ cattle scale service==online 5      # scale by label, which container has `-l service=online`
$ cattle scale name==nginx 5          # scale by container name, which container has `--name nginx`
$ cattle scale image==nginx:latest 5  # scale by image name
```
Scale down:
```
$ cattle scale service==online -5
```

~~~
Mutilple scale:
```
$ cattle scale service==online 5 service==offline -5
```
~~~

Filter: cattle scale complete compatible to swarm filter, just set ENV and labels to new container.

If the scale number < 0, will scale down containers on the node witch has `storage=ssd` label.
```
$ cattle scale -e constraint:storage==ssd -l app=scale-up-nginx service==online 5 
```

Http api.
```
```

