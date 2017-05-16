#!/bin/bash

docker pull $IMAGE
docker rm -f $NAME
docker run -d -p $PORT --name $NAME $IMAGE