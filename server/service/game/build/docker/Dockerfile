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

# build database migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1

# final image
FROM registry.gitlab.com/alienspaces/go-mud/postgres:12.7-alpine3.14

RUN apk update \
    && apk add curl \
    && apk add bash

COPY --from=builder /go/bin/go-mud-game-server /go/bin
COPY --from=builder /go/bin/go-mud-game-cli /go/bin
COPY --from=builder /go/bin/migrate /bin

# entrypoint
COPY --from=builder /go-mud/service/game/build/docker/entrypoint.sh .

# schemas
COPY --from=builder /go-mud/schema/docs ./schema/docs

# copy database migrations
COPY --from=builder /go-mud/service/game/migration /migration

EXPOSE 8084

# entrypoint
ENTRYPOINT [ "./entrypoint.sh" ]
