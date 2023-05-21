# Build image
ARG DOCKER_IMAGE_SERVICE_BASE_FULL_NAME
ARG DOCKER_IMAGE_SERVICE_SERVER_BASE_FULL_NAME
FROM $DOCKER_IMAGE_SERVICE_BASE_FULL_NAME AS builder

ENV APP_SERVER_HOME /pricing

WORKDIR $APP_SERVER_HOME
COPY . .

# Build server and CLI
WORKDIR $APP_SERVER_HOME/server/service/template
RUN go build -o /go/bin/pricing-template-server ./cmd/server
RUN go build -o /go/bin/pricing-template-cli ./cmd/cli

# Final image
FROM $DOCKER_IMAGE_SERVICE_SERVER_BASE_FULL_NAME

# Copy binaries
COPY --from=builder /go/bin/pricing-template-server /go/bin
COPY --from=builder /go/bin/pricing-template-cli /go/bin

# Copy entrypoint
COPY --from=builder /pricing/backend/service/template/build/docker/entrypoint.sh /go

# Copy JSON schemas
COPY --from=builder /pricing/backend/schema /go/schema

# Copy database migrations
COPY --from=builder /pricing/backend/service/template/migration /go/migration

# Entrypoint
ENTRYPOINT [ "./entrypoint.sh" ]
