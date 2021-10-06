#!/usr/bin/env bash

COMMAND=$1

echo "=> (entrypoint) Command ${COMMAND}"

if [ -z "$COMMAND" ]; then

    # run
    echo "=> (entrypoint) Executing run command"
    go-boilerplate-player-server

else

    # user command
    echo "=> (entrypoint) Executing user command $*"

    exec "$@"
fi
