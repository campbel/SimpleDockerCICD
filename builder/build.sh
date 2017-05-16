#!/bin/bash

git clone https://github.com/campbel/SimpleDockerCICD.git
cd SimpleDockerCICD/app
docker build -t campbel/app .
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker push campbel/app
