#!/bin/bash


##import /etc/profile
. /root/.bash_profile
. /etc/profile

##set
set -e

##work path
TARGET="/sdadata/item/proxyserver"
echo "cd ${TARGET}"
cd ${TARGET}


##compile
case "$1" in
    compile)
        echo "compile proxy(develop)..."
        ssh -t wangkebiao@192.168.1.39 "sudo /data/git/proxy/shell/proxy_dev_192.168.1.39.sh"
        if [ $? -ne 0 ]; then
            echo "proxy compile failed!!!"
            exit 1
        fi

        ##synchronize
        echo "proxy config is rsync..."
        rsync -av wangkebiao@192.168.1.39:/data/git/proxy/server/conf/*_dev* /sdadata/item/proxyserver/conf
        if [ $? -ne 0 ]; then
            echo "proxy rsync config failed!!!"
            exit 1
        fi

	      echo "proxy is rsync..."
        rsync -av wangkebiao@192.168.1.39:/data/git/proxy/server/proxy /sdadata/item/proxyserver
        if [ $? -ne 0 ]; then
            echo "proxy rsync proxy failed!!!"
            exit 1
        fi
	      ;;
esac



##start
proxy="/sdadata/item/proxyserver/proxy"
case "$1" in
    start)
        $proxy >> /sdadata/item/proxyserver/log/server.log 2>&1 &
        if [ $? -eq 0 ];then
            echo "proxy service start success!"
        else
            echo "proxy service start failed!"
            exit 1
        fi
        ;;
    stop)
        PID=`ps aux | grep "proxyserver/proxy" | grep -v grep | awk '{printf $2}'`
        if [ -n "$PID" ]; then
	          echo "kill -SIGINT ${PID}"
            sudo kill -SIGINT $PID
            if [ $? -eq 0 ];then
                echo "proxy service stop success!"
                /bin/sleep 10
            else
                echo "proxy service stop failed!"
                exit 1
            fi
        fi
        ;;
    restart)
        PID=`ps aux | grep "proxyserver/proxy" | grep -v grep | awk '{printf $2}'`
        if [ -n "$PID" ]; then
            echo "kill -SIGINT ${PID}"
	          sudo kill -SIGINT $PID
            if [ $? -eq 0 ];then
                echo "proxy service stop success!"
                /bin/sleep 10
                $proxy >> /sdadata/item/proxyserver/log/server.log 2>&1 &
                if [ $? -eq 0 ];then
                    echo "proxy service start success!"
                else
                    echo "proxy service start failed!"
                    exit 1
                fi
            else
                echo "proxy service stop failed!"
                exit 1
            fi
        fi
        ;;
esac


