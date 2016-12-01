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

## Feature of cattle
* Compatible docker api. 
* Simple.

## Container priority, min number

```
$ docker run -e PRIORITY=10 -e MIN_NUMBER=3 -l service=online --name nginx nginx:latest
$ docker run -e PRIORITY=1 -e MIN_NUMBER=1 -l service=offline --name nginx httpd:latest
```

## Scale up or down
By docker cli.
```
```
Http api.
```
```

