FROM alpine:3.5
RUN apk update && apk add git bash docker
ADD build.sh /build.sh
RUN mkdir workspace
WORKDIR workspace
CMD /bin/bash /build.sh