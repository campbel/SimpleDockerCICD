FROM jenkins:2.46.2-alpine

USER root
RUN apk update && apk add docker py-pip && pip install docker-compose

# install plugins
ADD plugins.txt /usr/share/jenkins/plugins.txt
RUN /usr/local/bin/plugins.sh /usr/share/jenkins/plugins.txt