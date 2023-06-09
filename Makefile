TAG = ${shell cat version}
DOCKER = docker
NAME = bookserver2:31000/tool/youtubedown

TEMPLATE = ./Dockerfile_tmp
TARGET = Dockerfile
TARGET_FILE = ./
GO_VERSION = 1.18.2


BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"



ARCH = ${shell arch}
ifeq (${ARCH},x86_64)
ARCH = amd64
BASECONTANA = ubuntu
else
ARCH = armv6l
BASECONTANA = ubuntu
endif

OPT = "--privileged"

create:
	@echo "create dockerfile"
	@echo "for ${NAME}:${TAG}"
	cat ${TEMPLATE} | sed -e "s|BASECONTANA|${BASECONTANA}|g" | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET_FILE}${TARGET}
build: create
	@echo "build"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET_FILE}${TARGET} ${TARGET_FILE}
up:
	${DOCKER} run -d --name=youtubedown -p 8080:8080 ${NAME}:${TAG} 
down: 
	${DOCKER} stop youtubedown
rm:
	${DOCKER} rm youtubedown
remove:
	${DOCKER} rmi ${NAME}:${TAG}
buildkit: create
	@echo "--- buildkit build --"
	${BUILD} --addr ${BUILD_ADDR_ARM} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${TARGET_FILE} --opt source=${TARGET_FILE}${TARGET}
