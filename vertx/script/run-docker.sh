#!/bin/bash

IMAGE_NAME="crate-vertx"
CONTAINER_NAME="crate-vertx-container"

echo "检查容器状态 并清理过期内容"

#echo "check if the container is running and stop it"
if [ $(docker ps -q -f name=$CONTAINER_NAME) ]; then
    docker stop $CONTAINER_NAME
fi

#echo "Check if the container exists and remove it"
if [ $(docker ps -a -q -f name=$CONTAINER_NAME) ]; then
    docker rm $CONTAINER_NAME
fi

#echo "Check if the image exists and remove it"
if [ $(docker images -q $IMAGE_NAME) ]; then
    docker rmi $IMAGE_NAME
fi

echo ""
echo "开始构建"
docker build -t $IMAGE_NAME .

echo ""
echo "运行"
docker run -p 8421:8421 --name $CONTAINER_NAME $IMAGE_NAME
