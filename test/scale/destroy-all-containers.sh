docker -H :4001 rm -f $(docker -H :4001 ps -a|grep ago|awk '{print $1}')
