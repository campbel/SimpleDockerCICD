FROM alpine:3.5
RUN apk update && apk add git bash docker
ADD deploy.sh /deploy.sh
RUN mkdir workspace
WORKDIR workspace
CMD /bin/bash /deploy.sh