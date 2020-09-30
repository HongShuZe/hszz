#!/bin/bash

function CleanContainer() {
    CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /education-peer.*/) {print $1}')
    if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
        echo "---没有可以删除的容器---"
    else
        docker rm -f $CONTAINER_IDS
    fi
}

function RemoveUnwantedImage() {
    DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /education-peer.*/) {print $3}')
    if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
        echo "---没有可以删除的镜像---"
    else
        docker rmi -f $DOCKER_IMAGE_IDS
    fi
}


CleanContainer

RemoveUnwantedImage

# 删除ca文件
rm -rf docker-compose-ca.yaml
rm -rf docker-compose-ca.yamlt