## API Refrence

### Scale
* URL: /scale
* METHOD: POST
* BODY:
```
{
  "Items": [
      {
         "Filters": ["service==online"],
         "Number": 3,
         "ENVs": ["constraint:storage==ssd", "affinity:app!=offline"],
         "Labels": {
            "app":"scale-up-nginx",
         }
      },
      {
         "Filters": ["service==offline"],
         "Number": -3,
         "ENVs": ["constraint:storage==ssd", "STOP_HOOK="www.iflytek.com/stop""],
         "Labels": {
            "app":"scale-down-nginx",
         }
      }
  ]
}
```
* Filters: define which container you want to scale.
* Number: how many container you want to scale.
* ENVs: set container enviroment, support swarm filters.
* Labels: set container label. If label already exist, overwrite it.

### Inform App Hook
* METHOD: POST
* BODY:
```
{
   "Action": "PRE_STOP" | "POST_STOP"
   "Containers": [
       {
             "Id": "8dfafdbc3a40",
             "Names":["/boring_feynman"],
             "Image": "ubuntu:latest",
             "ImageID": "d74508fb6632491cea586a1fd7d748dfc5274cd6fdfedee309ecdcbc2bf5cb82",
             "Command": "echo 1",
             "Created": 1367854155,
             "State": "Exited",
             "Status": "Exit 0",
             "Ports": [{"PrivatePort": 2222, "PublicPort": 3333, "Type": "tcp"}],
             "Labels": {
                     "com.example.vendor": "Acme",
                     "com.example.license": "GPL",
                     "com.example.version": "1.0"
             },
             "SizeRw": 12288,
             "SizeRootFs": 0,
             "Node": {
                 "Id": "ODAI:IC6Q:MSBL:TPB5:HIEE:6IKC:VCAM:QRNH:PRGX:ERZT:OK46:PMFX",
                 "Ip": "0.0.0.0",
                 "Addr": "http://0.0.0.0:4243",
                 "Name": "vagrant-ubuntu-saucy-64"
             },
             "HostConfig": {
                     "NetworkMode": "default"
             },
             "NetworkSettings": {
                     "Networks": {
                             "bridge": {
                                      "IPAMConfig": null,
                                      "Links": null,
                                      "Aliases": null,
                                      "NetworkID": "7ea29fc1412292a2d7bba362f9253545fecdfa8ce9a6e37dd10ba8bee7129812",
                                      "EndpointID": "2cdc4edb1ded3631c81f57966563e5c8525b81121bb3706a9a9a3ae102711f3f",
                                      "Gateway": "172.17.0.1",
                                      "IPAddress": "172.17.0.2",
                                      "IPPrefixLen": 16,
                                      "IPv6Gateway": "",
                                      "GlobalIPv6Address": "",
                                      "GlobalIPv6PrefixLen": 0,
                                      "MacAddress": "02:42:ac:11:00:02"
                              }
                     }
             },
             "Mounts": [
                     {
                              "Name": "fac362...80535",
                              "Source": "/data",
                              "Destination": "/data",
                              "Driver": "local",
                              "Mode": "ro,Z",
                              "RW": false,
                              "Propagation": ""
                     }
             ]
        }, 
   ]
}
```
* Action: PRE_STOP before stop the containers, POST_STOP after stop containers
* Node: contain the container Host imformation
