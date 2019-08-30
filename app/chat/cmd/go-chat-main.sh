#!/bin/bash

### BEGIN INIT INFO
# Provides:          go-chat-main
# Required-Start:    $local_fs $network
# Required-Stop:     $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: go-chat-main service
# Description:       go-chat-main service
### END INIT INFO

# service's name
srv_name="go-chat-main"
# workdir
base_path=$(cd `dirname $0`; pwd)
# where the bin file is, the name depends on you compile
bin_path="$base_path/main-darwin"
# create log file by date `date '+%Y-%m-%d'`
log_file="$base_path/logs/`date '+%Y-%m-%d'`.log"

start(){

    echo "Starting $srv_name ..."

    # using nohup daemonize
    nohup $bin_path >> $log_file 2>&1 &

    echo "Runing with command $bin_path >> $log_file 2>&1 &"

    # if command 'nohup' not return anything
    # it means service start faild
    if [ "$?" != 0 ] ; then
        echo " failed"
        exit 1
    fi

    echo 'OK'
}


stop(){

    echo "Stoping $srv_name ..."

    # find pid
    pid=`ps aux|grep "$bin_path"|grep -v grep|awk '{print $2}'`

    if [ ! -n "$pid" ] ; then
        echo "warning, no pid found - are u sure this service is running ?"
        exit 1
    fi

    # kill the process
    kill -9 $pid

    echo 'OK'
}


restart(){
    stop
    start
}


case $1 in

    start)
        start
    ;;

    stop)
        stop
    ;;

    restart)
        restart
    ;;

    *)

    echo "Usage: $0 {start|stop|restart}"
    exit 1
    ;;

esac