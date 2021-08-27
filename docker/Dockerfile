FROM golang:bullseye

ARG DOCKER_UID
RUN mkdir /.cache && chown -R ${DOCKER_UID}:${DOCKER_UID} /.cache && chmod -R 0777 /.cache

WORKDIR /go/src/app

ENTRYPOINT ["go", "run", "main.go"]