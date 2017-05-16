#!/bin/bash

git clone $REPOSITORY
cd $DOCKERFILE
docker build -t $IMAGE .
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker push $IMAGE
