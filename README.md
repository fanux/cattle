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

## Namespace, service, app 
* Namespace is a collection of containers or nodes. `docker run -l namespace=swarm swarm:latet`
* Service is a collection of different kinds of containers. For example, you can define `nginx,mysql,php` container as a same service :`docker run -l service=web-service ...`
* App is a collection of same kind of containers. For example, you run 5 replication of nginx, all set the `app=nginx` lalel.

## Container priority, min number

```
$ docker run -e PRIORITY=10 -e MIN_NUMBER=3 -l service=online --name nginx nginx:latest
$ docker run -e PRIORITY=1 -e MIN_NUMBER=1 -l service=offline --name nginx httpd:latest
```

## Scale up or down
Cli

Scale up: suggest use docker compose service as a scale unit.
```
$ cattle scale [[env|label] filter number]
```
```
$ cattle scale -f service==online -n 5      # scale by label, which container has `-l service=online`
$ cattle scale -f name==nginx -n 5          # scale by container name, which container has `--name nginx`
~~$ cattle scale -f image==nginx:latest -n 5  # scale by image name~~
```
The same app will not scale up again, judged by lable `-l app=xxx`:
```
  php:2 total                                            php:2+3=5 total                                    
  redis:1 total                                          redis:1+3=4 total                                  
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  | service=web |                                        | service=web | | service=web | | service=web |    
  | app=php     |                                        | app=php     | | app=php     | | app=redis   |    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  | service=web |  cattle scale -f service=web -n 3      | service=web | | service=web | | service=web |    
  | app=php     |====================================>   | app=php     | | app=php     | | app=redis   |    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  | service=web |                                        | service=web | | service=web | | service=web |    
  | app=redis   |                                        | app=php     | | app=redis   | | app=redis   |    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
```

Scale down:
```
$ cattle scale service==online -5
```
Filter: cattle scale complete compatible to swarm filter, just set ENV and labels to new container.

If the scale number < 0, will scale down containers on the node witch has `storage=ssd` label.
```
$ cattle scale -e constraint:storage==ssd -l app=scale-up-nginx -f service==online -n 5 
```
`--force` this flag will stop those containers witch priority is below then scale up service when resource is not enough.

Http api.
```
```

