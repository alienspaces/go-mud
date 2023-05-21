#!/usr/bin/env bash

COMMAND=$1

echo "=> (entrypoint) APP_ENV ${APP_ENV}"

# Source environment
ENV_FILE=".env.${APP_ENV}"
if [ -f "$ENV_FILE" ]; then
    echo "=> (entrypoint) Sourcing .env"
    # Copy the environment file to .env so application binaries also detect it
    cp "$ENV_FILE" .env
    source .env
fi

echo "=> (entrypoint) env"
env | grep APP | sort

echo "=> (entrypoint) Command ${COMMAND}"

if [ -z "$COMMAND" ]; then

    # NOTE:
    # To support deploying to a QA environment from Slack the deployment job has to be
    # rolled into a single Gitlab job as opposed to the typical test, build, migrate
    # then deploy jobs. The QA deployment job is run within the MSTS dind container
    # which does not contain the Go build chain so does not contain the necessary tools
    # for executing database migrations. Normally the migration job would be run inside
    # an Alpine/Go container. The QA environment database is also refreshed with each
    # deployment so we will migrate down completely, then migrate up, apply database
    # grants and load test data.

    if [ "${APP_ENV}" == "qa" ]; then

        URL="postgres://$APP_SERVER_DB_OWNER_USER:$APP_SERVER_DB_OWNER_PASSWORD@$APP_SERVER_DB_HOST:$APP_SERVER_DB_PORT/$APP_SERVER_DB_NAME?sslmode=disable"

        echo "=> (entrypoint) Migration URL: ${URL}"
        echo "=> (entrypoint) Migration PATH: ./migration"

        echo "=> (entrypoint) Migrating down"
        # shellcheck disable=SC2012
        MIGRATIONCOUNT=$(ls ./migration | wc -l);
        migrate -verbose -path "./migration" -database "$URL" down "$MIGRATIONCOUNT"

        echo "=> (entrypoint) Migrating up"
        migrate -verbose -path "./migration" -database "$URL" up

        echo "=> (entrypoint) Migration applying database grants to $APP_SERVER_DB_USER"

        export PGPASSWORD=$APP_SERVER_DB_OWNER_PASSWORD
        psql --host="$APP_SERVER_DB_HOST" \
            --port="$APP_SERVER_DB_PORT" \
            --username="$APP_SERVER_DB_OWNER_USER" "$APP_SERVER_DB_NAME" \
            --command "GRANT SELECT, INSERT, UPDATE, TRUNCATE, DELETE ON ALL TABLES IN SCHEMA \"public\" TO \"$APP_SERVER_DB_USER\";"

        echo "=> (entrypoint) Migration loading test data"
        pricing-template-cli db-load-test-data
    fi

    # Run server
    echo "=> (entrypoint) Executing run command"
    pricing-template-server

else

    # Run user command
    echo "=> (entrypoint) Executing user command $*"

    exec "$@"
fi
