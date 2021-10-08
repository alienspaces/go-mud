# build image
FROM registry.gitlab.com/alienspaces/go-mud/golang:1.17.2-alpine3.14 as builder

# enable https:// repositories
RUN apk update \
    && apk add curl \
    && apk add bash \
    && apk add git

# app home
ENV APP_SERVER_HOME /go-mud

# copy application
WORKDIR $APP_SERVER_HOME
COPY . .

# build API server
WORKDIR $APP_SERVER_HOME/service/game
RUN go build -o /go/bin/go-mud-game-server ./cmd/server
RUN go build -o /go/bin/go-mud-game-cli ./cmd/cli

# final image
FROM registry.gitlab.com/alienspaces/go-mud/golang:1.17.2-alpine3.14

RUN apk update \
    && apk add curl \
    && apk add bash

COPY --from=builder /go/bin/go-mud-game-server /go/bin
COPY --from=builder /go/bin/go-mud-game-cli /go/bin
COPY --from=builder /go-mud/service/game/build/docker/entrypoint.sh .
COPY --from=builder /go-mud/schema/docs ./schema/docs

RUN pwd
RUN ls -la

# entrypoint
ENTRYPOINT [ "./entrypoint.sh" ]
