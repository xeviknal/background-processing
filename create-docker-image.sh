DOCKER_NAME=xeviknal/background-processor
DOCKER_VERSION=v0.1
DOCKER_TAG=${DOCKER_NAME}:${DOCKER_VERSION}

docker build -t ${DOCKER_TAG} .
