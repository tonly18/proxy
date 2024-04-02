PROJECT_NAME:="proxy-server"
VERSION:=v1

.PHONY: build start restart stop clean image run

build:
	bash build.sh compile ${PROJECT_NAME}


start:
	bash build.sh start ${PROJECT_NAME}


restart:
	bash build.sh restart ${PROJECT_NAME}


stop:
	bash build.sh stop ${PROJECT_NAME}


clean:
	bash build.sh clean ${PROJECT_NAME}


image:
	docker build -t ${PROJECT_NAME}:${VERSION} .

run:
	docker run -itd -p 8080:8080 ${PROJECT_NAME}:${VERSION}
