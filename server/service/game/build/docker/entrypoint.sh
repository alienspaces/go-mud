#!/usr/bin/env bash

# Retry a command a number of times with an
# increasing wait time between each attempt.
function retry_cmd {
  local n=1
  local max=5
  local delay=5
  local delay_inc=5

  echo "=> Command $*"

  while true; do
    "$@" && break || {
      if [[ $n -lt $max ]]; then
        ((n++))
        echo "=> Command failed. Trying again in $delay seconds. Attempt $n/$max:"
        sleep $delay;
        delay=$[$delay+$delay_inc]
      else
        echo "=> Command failed after $n attempts, exiting.." >&2
        exit 1
      fi
    }
  done
}

COMMAND=$1

echo "=> (entrypoint) Command ${COMMAND}"

echo "=> (entrypoint) PWD ${PWD}"

ls -la

if [ -z "$COMMAND" ]; then

    # postgres
    echo "=> (entrypoint) Starting Postgres"

    nohup ./docker-entrypoint.sh postgres &

    # extensions
    echo "=> (entrypoint) Creating extension pgcrypto"
    retry_cmd psql --host="$APP_SERVER_DB_HOST" --port="$APP_SERVER_DB_PORT" --username="$APP_SERVER_DB_USER" --command="CREATE EXTENSION pgcrypto;" "$APP_SERVER_DB_NAME"

    # migrate
    URL="postgres://$APP_SERVER_DB_USER:$APP_SERVER_DB_PASSWORD@$APP_SERVER_DB_HOST:$APP_SERVER_DB_PORT/$APP_SERVER_DB_NAME?sslmode=disable"

    echo "=> (entrypoint) Running migrations ${URL}"
    migrate -verbose -path "./migration" -database "$URL" -- up

    # load test data
    echo "=> (entrypoint) Loading test data"
    go-mud-game-cli load-test-data || exit $?

    # run server
    echo "=> (entrypoint) Executing run command"
    go-mud-game-server

else

    # user command
    echo "=> (entrypoint) Executing user command $*"

    exec "$@"
fi
