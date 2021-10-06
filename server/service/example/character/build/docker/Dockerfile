# build image
FROM golang:1.16.4-alpine3.13 as builder

# enable https:// repositories
RUN apk update \
    && apk add curl \
    && apk add bash \
    && apk add git

# app home
ENV APP_SERVER_HOME /go-boilerplate

# copy application
WORKDIR $APP_SERVER_HOME
COPY . .

# build API server
WORKDIR $APP_SERVER_HOME/service/character
RUN go build -o /go/bin/go-boilerplate-character-server ./cmd/server
RUN go build -o /go/bin/go-boilerplate-character-cli ./cmd/cli

# final image
FROM golang:1.16.4-alpine3.13

RUN apk update \
    && apk add curl \
    && apk add bash

COPY --from=builder /go/bin/go-boilerplate-character-server /go/bin
COPY --from=builder /go/bin/go-boilerplate-character-cli /go/bin
COPY --from=builder /go-boilerplate/service/character/build/docker/entrypoint.sh .
COPY --from=builder /go-boilerplate/schema/docs ./schema/docs

RUN pwd
RUN ls -la

# entrypoint
ENTRYPOINT [ "./entrypoint.sh" ]
