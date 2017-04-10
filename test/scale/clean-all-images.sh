docker -H :4001 rmi -f $(docker -H :4001 images|grep ago|awk '{print $1":"$2}')
