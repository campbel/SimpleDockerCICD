FROM jenkins:2.46.2-alpine

USER root

# install docker
RUN apk add --no-cache \
		ca-certificates

ENV DOCKER_CHANNEL edge
ENV DOCKER_VERSION 17.05.0-ce

RUN set -ex; \
	apk add --no-cache --virtual .fetch-deps \
		curl \
		tar \
	; \
	curl -fL -o docker.tgz "https://download.docker.com/linux/static/${DOCKER_CHANNEL}/x86_64/docker-${DOCKER_VERSION}.tgz"; \
	tar --extract \
		--file docker.tgz \
		--strip-components 1 \
		--directory /usr/local/bin/ \
	; \
	rm docker.tgz; \
	apk del .fetch-deps; \
	dockerd -v; \
	docker -v

# install docker-compose
RUN apk update && apk add py-pip && pip install docker-compose

# install jenkins plugins
ADD plugins.txt /usr/share/jenkins/plugins.txt
RUN /usr/local/bin/plugins.sh /usr/share/jenkins/plugins.txt