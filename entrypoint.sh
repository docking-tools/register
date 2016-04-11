#!  /bin/sh
set -e

if [ "$1" = 'dev' ]; then
    shift
    register -c /register/config.json
        --hostip="$HOST_IP" \
        --register="$REGISTER_URL" \
        --docker="$DOCKER_URL" \
        "$@"
else
    exec "$@"
fi