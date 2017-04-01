[![Build Status](http://10.1.86.50:8000/api/badges/fanux/cattle/status.svg)](http://10.1.86.50:8000/fanux/cattle)
## Why cattle  
1. Swarm does not support scale command or API. Cattle support scale service with filters.
2. Scale up is easy, but when some high priority service want to seize low priority resources, how to decide which service to scale down?
3. Stop container after inform app.

Cattle solve those problems. 

## Special labels: Namespace, service, app
* Namespace is a collection of containers or nodes. `docker run -l namespace=swarm swarm:latest`
* Service is a collection of different kinds of containers. For example, you can define `nginx,mysql,php` container as the same service :`docker run -l service=web-service ...`
* App is a collection of the same kind of containers. For example, you run 5 replication of nginx, all set the `app=nginx` label.

## Container priority, min number

**MIN_NUMBER is going to be discard, if you don't want your app to be seized, just set a priority=0**
```
$ docker run -e PRIORITY=10 -e MIN_NUMBER=3 -l service=online --name nginx nginx:latest
$ docker run -e PRIORITY=1 -e MIN_NUMBER=1 -l service=offline --name nginx httpd:latest
```

## Scale up or down
Scale up: suggest use docker compose service as a scale unit.
```
$ cattle scale [[env|label] filter number]
```
```
$ cattle scale -f service==online -n 5      # scale by service, which container has `-l service=online`
$ cattle scale -f name==nginx -n 5          # scale by container name, which container has `--name nginx`
$ cattle scale -f foo==bar -n 5             # scale by the label, just select a container has `foo=bar` label, and scale it(not all the containers!).
```
The same app will not scale up again, judged by label `-l app=xxx`:
```
  php:2 total                                            php:2+3=5 total                                    
  redis:1 total                                          redis:1+3=4 total                                  
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  | service=web |                                        | service=web | | service=web | | service=web |    
  | app=php     |                                        | app=php     | | app=php     | | app=redis   |    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  | service=web | cattle scale -f service==web -n 3      | service=web | | service=web | | service=web |    
  | app=php     |====================================>   | app=php     | | app=php     | | app=redis   |    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
  | service=web |                                        | service=web | | service=web | | service=web |    
  | app=redis   |                                        | app=php     | | app=redis   | | app=redis   |    
  +-------------+                                        +-------------+ +-------------+ +-------------+    
```

Scale down:
```
$ cattle scale -f service==online -n -5
```
Filter: cattle scale complete compatible to swarm filter, just set ENV and labels to new container.

If the scale number < 0, will scale down containers on the node witch has `storage=ssd` label.
```
$ cattle scale -e constraint:storage==ssd -l app=scale-up-nginx -f service==online -n 5 
```


If the container set the ENV MIN_NUMBER=x, cattle will guarantee has x containers left after scale down. 
```
 php:5 total MIN_NUMBER=3                                                            php:3 left
 redis:4 total MIN_NUMBER=1                                                          redis:1 left
 +-------------+ +-------------+ +-------------+                                     +-------------+ +-------------+
 | service=web | | service=web | | service=web |                                     | service=web | | service=web |
 | app=php     | | app=php     | | app=redis   |                                     | app=php     | | app=redis   |
 +-------------+ +-------------+ +-------------+                                     +-------------+ +-------------+
 +-------------+ +-------------+ +-------------+                                     +-------------+
 | service=web | | service=web | | service=web | cattle scale -f service==web -n -3  | service=web |
 | app=php     | | app=php     | | app=redis   | ==================================> | app=php     |
 +-------------+ +-------------+ +-------------+ php: 5 - 3 < MIN_NUMBER=3           +-------------+
 +-------------+ +-------------+ +-------------+ redis: 4 - 3 = MIN_NUMBER=1         +-------------+
 | service=web | | service=web | | service=web |                                     | service=web |
 | app=php     | | app=redis   | | app=redis   |                                     | app=php     |
 +-------------+ +-------------+ +-------------+                                     +-------------+
```

### Scale with constraint
scale up to nodes whatever you want, scale down containers on nodes you assign.

scale up nginx to node without GPU:
```
$ cattle scale -f app==nginx -e constraint:GPU!=true -n 5
```

only scale down nginx on centos7 node:
```
$ cattle scale -f app==nginx -e constraint:os==centos7 -n -5
```

## Resource Seize
Resource Seize is complex, a Resource Seize scale task must has those argument:

* constraint : what node resource you want to seize.
* inaffinity : what container resource you want to seze.
* applots    : run how many container on one node. 
* priority   : high priority can seize low priority resource.

If the free node is not enough( free-node-num * applots < need-scale-up-number), will tick the seize.
What is a free node? The node don't have the inaffinity containers.

User `affinity:xxx!=xxx` will stop those containers witch priority is below then scale up service.

Suggest there are four nodes with `GPU=true` label in our cluster. There are 2 services: online and offline. The MIN_NUMBER of offline is 1, onlie has higher priority 1 and offline has lower priority 9.
```
 node: GPU=true            node: GPU=true            node: GPU=true            node: GPU=true
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 | | service=online |    | | | service=offline|    | | | service=offline|    | | | service=offline|    |
 | | priority=1     |    | | | priority=9     |    | | | priority=9     |    | | | priority=9     |    |
 | |                |    | | | MIN_NUMBER=1   |    | | | MIN_NUMBER=1   |    | | | MIN_NUMBER=1   |    |
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
                                                    | cattle scale -e constraint:GPU==true       \     # I need GPU nodes
                                                    |              -e applots=1                  \     # One node one container
                                                    |              -e affinity:service!=offline  \     # I seize the offline resource
                                                    |              -f service==online  -n 3            # Scale up 3 online service instances
                                                    V
 node: GPU=true            node: GPU=true            node: GPU=true            node: GPU=true
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+   Want scale up 3 online service instances, but only 2 successed, because must
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |   ensure the MIN_NUMBER of offline service left.
 | | service=online |    | | | service=online |    | | | service=online |    | | | service=offline|    |
 | | priority=1     |    | | | priority=1     |    | | | priority=1     |    | | | priority=9     |    |
 | |                |    | | |                |    | | |                |    | | | MIN_NUMBER=1   |    |
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
                                                    | cattle scale -e constraint:GPU==true \           # Want scale up offline service
                                                    |              -e applots=1                  \     # One node one container
                                                    |              -e affinity:service!=online \
                                                    |              -f service==offline -n 2
                                                    V
 node: GPU=true            node: GPU=true            node: GPU=true            node: GPU=true
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+ Seize resource failed. Because the offline
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    | priority is lower. Offline service must wait
 | | service=online |    | | | service=online |    | | | service=online |    | | | service=offline|    | online service initiative release its resource.
 | | priority=1     |    | | | priority=1     |    | | | priority=1     |    | | | priority=9     |    |
 | |                |    | | |                |    | | |                |    | | | MIN_NUMBER=1   |    |
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
                                                    | cattle scale -f service==online -n -2
                                                    V
 node: GPU=true            node: GPU=true            node: GPU=true            node: GPU=true
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+ Online service initiative release its resource.
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 | | service=online |    | | |                |    | | |                |    | | | service=offline|    |
 | | priority=1     |    | | |                |    | | |                |    | | | priority=9     |    |
 | |                |    | | |                |    | | |                |    | | | MIN_NUMBER=1   |    |
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
                                                    | cattle scale -f service==offline -n 2
                                                    V
 node: GPU=true            node: GPU=true            node: GPU=true            node: GPU=true
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 | | service=online |    | | | service=offline|    | | | service=offline|    | | | service=offline|    |
 | | priority=1     |    | | | priority=9     |    | | | priority=9     |    | | | priority=9     |    |
 | |                |    | | | MIN_NUMBER=1   |    | | | MIN_NUMBER=1   |    | | | MIN_NUMBER=1   |    |
 | +----------------+    | | +----------------+    | | +----------------+    | | +----------------+    |
 +-----------------------+ +-----------------------+ +-----------------------+ +-----------------------+
```

## Inform App Hook 
Before stop a container, must inform it to do some clean work.
```
-e STOP_HOOK="www.iflytek.com/stop" \
-e WAIT_TIME=60s
```

```
     www.iflytek.com/stop           cattle
             |        PRE_STOP        |
             |<-----------------------|
             |     container info     | sleep 60s
             |                        |
             |       POST_STOP        | stop the container
             |<-----------------------|
             |    container info      |
             V                        V
```

### applots support
```
$ docker run -l app=foo -e "applots=3" foo:latest
```
One host run less then 3 containers which has `app=foo` label. (`app` is a special label)

### create containers with replication
```
$ docker run -e "replica=3" foo:latest
```
This command will create 3 containers using `foo:latest` image.

### support task types
Scale up container has two types task currently. Create container or start a stoped container.

By default cattle create new container. If you don't want to create new containers, using -e TASK_TYPE=start

```
$ cattle scale -f key==value -e TASK_TYPE=start -n 5
$ cattle scale -f key==value -e TASK_TYPE=create -n 5
```

Scale down has many task types as well, cattle destroy containers by default, if you want stop container:
```
$ cattle scale -f key==value -e TASK_TYPE=stop -n -5
$ cattle scale -f key==value -e TASK_TYPE=destroy -n -5
```

### TimeSlice
Container will stop or remove after TimeSlice.
```
$ cattle scale -f key==value -e TIMESLICE=2h -n 5 
```
This is usful for prevent high priority app don't release resource.

## task add and action
Add command not do tasks immediately, send request to manager when excute action command!
this is useful for rolling update. manager will scale up and down alternately
```
$ cattle add -f app==foov2 -n 5
$ cattle add -f app==foov1 -n -5
```
```
$ cattle action
```

## scale file
Scale file is a yaml config file. Touch a file named `scale-file.yml`:

```
version: v1
items:
    scaleUpFoov2:
        filters:
            - "app==foov2"
        number: 5
        labels: 
            - "key=value"
        envs:
            - "constraint:GPU==true"
        
    scaleDownFoov1:
        ...
```

At scale: default file name is `scale-file.yml`
```
$ cattle -f scale-file.yml scale
```
## Let's run scale!!!
Using docker cli to scale! 

know that `scale` is not a real image name, if cattle get a image named `scale`, just run scale action instead of create a real container! 

so we can using docker cli, and need not to install cattle cli. It convenient.
```
docker run scale -f ...
```
