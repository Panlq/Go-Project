#!/bin/sh
...
# allow the container to be started with `--user`
# if [ "$1" = 'hello' -a "$(id -u)" = '0' ]; then
#     exec gosu app "$@"
# fi

exec "$@"