SERVICE_NAME:=proxy-server
BRANCH:=master

.PHONY: build start restart stop clean

default: build

build:
	bash build.sh compile ${SERVICE_NAME} ${BRANCH}


start:
	bash build.sh start ${SERVICE_NAME}


restart:
	bash build.sh restart ${SERVICE_NAME}


stop:
	bash build.sh stop ${SERVICE_NAME}


clean:
	bash build.sh clean ${SERVICE_NAME}
