#!/bin/bash


##import /etc/profile
. /root/.bash_profile
. /etc/profile

##set
set -e

##PATH
APP_WORKPATH="/data/service/proxy"
RUNTIME_LOG="${APP_WORKPATH}/log/runtime.log"

##variable
APP_NAME=$2

##build start restart stop clean
case "$1" in
    compile)
        ##variable
        BUILD_VERSION=$(git log -1 --oneline)
        GIT_REVISION=$(git rev-parse --short HEAD)
        GIT_BRANCH=$(git name-rev --name-only HEAD)
        BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S")
        GO_VERSION=$(go version)

        ##go build
        go mod tidy
        go mod download
        go build -a -ldflags '-extldflags "-static"' -ldflags " \
        	-X 'main.AppName=${APP_NAME}'	\
        	-X 'main.AppVersion=${APP_VERSION}' \
        	-X 'main.BuildVersion=${BUILD_VERSION//\'/_}' \
        	-X 'main.BuildTime=${BUILD_TIME}' \
        	-X 'main.GitRevision=${GIT_REVISION}' \
          -X 'main.GitBranch=${GIT_BRANCH}' \
        	-X 'main.GoVersion=${GO_VERSION}' \
        	" -o $APP_NAME
        if [ $? -ne 0 ]; then
            exit 1
        fi

        ##version
        echo "${BUILD_TIME} ${GO_VERSION}" >> "${APP_WORKPATH}/version"
        if [ $? -eq 0 ]; then
            echo "${APP_NAME} go build success!"
        else
            echo "${APP_NAME} go build failed!"
            exit 1
        fi
	      ;;
	  start)
	      ##work path
        echo "cd ${APP_WORKPATH}"
        cd ${APP_WORKPATH}

        ##start
        SERVICE_CMD="${APP_WORKPATH}/${APP_NAME}"
        ${SERVICE_CMD} >> ${RUNTIME_LOG} 2>&1 &
        if [ $? -eq 0 ];then
            /bin/sleep 1
            echo "${APP_NAME} service start success!"
        else
            echo "${APP_NAME} service start failed!"
            exit 1
        fi
        ;;
    stop)
        ##work path
        echo "cd ${APP_WORKPATH}"
        cd ${APP_WORKPATH}

        ##stop
        PID=$(ps x | grep $APP_NAME | grep -v build.sh | grep -v grep | awk '{print $1}')
        if [ -n "$PID" ]; then
            echo "kill -9 ${PID}"
            sudo kill -9 $PID
            if [ $? -eq 0 ];then
                /bin/sleep 5
                echo "${APP_NAME} service stop success!"
            else
                echo "${APP_NAME} service stop failed!"
                exit 1
            fi
        else
            echo "${APP_NAME} service stop error"
            exit 1
        fi
        ;;
    restart)
        ##work path
        echo "cd ${APP_WORKPATH}"
        cd ${APP_WORKPATH}

        ##restart
        make stop
        if [ $? -eq 0 ];then
            /bin/sleep 1
            make start
            if [ $? -ne 0 ];then
              exit 1
            fi
        fi
        ;;
    clean)
        ##work path
        echo "cd ${APP_WORKPATH}"
        cd ${APP_WORKPATH}

        ##delete app
        sudo rm -rf ${APP_NAME}
         if [ $? -eq 0 ];then
            echo "${APP_NAME} clean service success!"
         else
            echo "${APP_NAME} clean service failed!"
            exit 1
        fi
        ;;
esac
