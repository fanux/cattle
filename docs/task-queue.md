### The Ring task queue

In this case, show us how cattle make a rolling update.

There are four old containers and a new container in the cluster.

```
┌─────────────┐   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ 
│ foo:v1      │   │ foo:v1      │  │ foo:v1      │  │ foo:v1      │  │ foo:v2      │
└─────────────┘   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```
Now we want update those container with `foo:v2`

```
$ cattle -f foo==v1 -n -4  -f foo==v2 -n 4       //current command line not support multiple scale items. 
```
So cattle product 8 tasks, four tasks to stop v1 container, and four tasks to start v2 container. 

Cattle do stop task and start task alternately.

The Ring task queue
```
               ┌─────────────┐  ┌─────────────┐               
       +------>│ start v2    │->│ stop v1     │------+       
       |       └─────────────┘  └─────────────┘      |        
       |                                             V
┌─────────────┐                                ┌─────-───────┐
│ stop v1     │                                │ start v2    │
└─────────────┘                                └─────────────┘
       ^                     Ring                    | 
       |                                             V
┌─────────────┐                                ┌─────-───────┐
│ start v2    │                                │ stop v1     │
└─────────────┘                                └─────────────┘
       ^                                             |
       |        ┌─────────────┐  ┌─────────────┐     |         
       +--------│ stop v1     │<-│ start v2    │<----+        
                └─────────────┘  └─────────────┘               
```
Retry, if a task execute success, will remove from Ring queue, If failed the task retry count--, if retry count is 0,
remove from Ring queue, and report this event!

### Advantage
* Support Rolling update.
* We can't done start task then start do stop task, because !affinity may cause create container failed, if old container not stop. 
* We can't done stop task then start do start task, will not suport rolling update.
* If Seize resource, need consider the WAIT time, so each stop task need do in a thread.
